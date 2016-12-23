package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielkrainas/gobag/api/errcode"
	"github.com/danielkrainas/gobag/context"
	"github.com/danielkrainas/gobag/decouple/cqrs"
	"github.com/gorilla/mux"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/auth"
	"github.com/danielkrainas/tinkersnest/configuration"
	"github.com/danielkrainas/tinkersnest/queries"
	"github.com/danielkrainas/tinkersnest/storage"
)

type dispatchFunc func(ctx context.Context, r *http.Request) http.Handler

type App struct {
	context.Context // TODO: does this need to be a context?

	config *configuration.Config

	router *mux.Router
}

func (app *App) Value(key interface{}) interface{} {
	if ks, ok := key.(string); ok && ks == "server.app" {
		return app
	}

	return app.Context.Value(key)
}

func getApp(ctx context.Context) *App {
	if app, ok := ctx.Value("server.app").(*App); ok {
		return app
	}

	return nil
}

type appRequestContext struct {
	context.Context

	URLBuilder *v1.URLBuilder
}

func (arc *appRequestContext) Value(key interface{}) interface{} {
	switch key {
	case "url.builder":
		return arc.URLBuilder
	}

	return arc.Context.Value(key)
}

func getURLBuilder(ctx context.Context) *v1.URLBuilder {
	if ub, ok := ctx.Value("url.builder").(*v1.URLBuilder); ok {
		return ub
	}

	return nil
}

func NewApp(ctx context.Context, config *configuration.Config) (*App, error) {
	app := &App{
		Context: ctx,
		config:  config,
		router:  v1.RouterWithPrefix(""),
	}

	app.register(v1.RouteNameBase, func(ctx context.Context, r *http.Request) http.Handler {
		return http.HandlerFunc(apiBase)
	})

	app.register(v1.RouteNameBlog, blogListDispatcher)
	app.register(v1.RouteNamePostByName, postByNameDispatcher)
	app.register(v1.RouteNameUserRegistry, userRegistryDispatcher)
	app.register(v1.RouteNameUserByName, userByNameDispatcher)
	app.register(v1.RouteNameAuth, authDispatcher)
	return app, nil
}

func (app *App) authorizeUser(ctx *appRequestContext, r *http.Request) error {
	route := mux.CurrentRoute(r)
	routeName := route.GetName()
	bearer := r.Header.Get("Authorization")
	authParts := strings.Split(bearer, ":")
	bearer = strings.TrimSpace(authParts[len(authParts)-1])
	if bearer == "" {
		_, hasClaim := ctx.Value("claim").(*v1.Claim)
		if hasClaim && r.Method == http.MethodPost {
			return nil
		} else if routeName != v1.RouteNameAuth {
			return errors.New("invalid bearer token")
		}
	}

	userName, err := auth.VerifyBearerToken(bearer)
	if err != nil {
		return err
	}

	user, err := cqrs.DispatchQuery(ctx, queries.FindUser{userName})
	if err != nil {
		return err
	}

	ctx.Context = context.WithValue(ctx.Context, "user", user)
	ctx.Context = acontext.WithLogger(ctx.Context, acontext.GetLoggerWithField(ctx.Context, "user.name", userName))
	return nil
}

func (app *App) dispatcher(dispatch dispatchFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.context(w, r)
		ctx.Context = acontext.WithErrors(ctx.Context, make(errcode.Errors, 0))

		if err := preloadClaim(ctx, r); err != nil {
			acontext.GetLogger(ctx).Error(err)
			ctx.Context = acontext.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		} else if err := app.authorizeUser(ctx, r); err != nil {
			acontext.GetLogger(ctx).Error(err)
			ctx.Context = acontext.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		} else {
			dispatch(ctx, r).ServeHTTP(w, r)
		}

		if errors := acontext.GetErrors(ctx); errors.Len() > 0 {
			if err := errcode.ServeJSON(w, errors); err != nil {
				acontext.GetLogger(ctx).Errorf("error serving error json: %v (from %s)", err, errors)
			}

			app.logError(ctx, errors)
		}
	})
}

func (app *App) logError(ctx context.Context, errors errcode.Errors) {
	for _, err := range errors {
		var lctx context.Context

		switch err.(type) {
		case errcode.Error:
			e, _ := err.(errcode.Error)
			lctx = acontext.WithValue(ctx, "err.code", e.Code)
			lctx = acontext.WithValue(lctx, "err.message", e.Code.Message())
			lctx = acontext.WithValue(lctx, "err.detail", e.Detail)
		case errcode.ErrorCode:
			e, _ := err.(errcode.ErrorCode)
			lctx = acontext.WithValue(ctx, "err.code", e)
			lctx = acontext.WithValue(lctx, "err.message", e.Message())
		default:
			// normal "error"
			lctx = acontext.WithValue(ctx, "err.code", errcode.ErrorCodeUnknown)
			lctx = acontext.WithValue(lctx, "err.message", err.Error())
		}

		lctx = acontext.WithLogger(ctx, acontext.GetLogger(lctx,
			"err.code",
			"err.message",
			"err.detail"))

		acontext.GetResponseLogger(lctx).Errorf("response completed with error")
	}
}

func (app *App) context(w http.ResponseWriter, r *http.Request) *appRequestContext {
	ctx := acontext.DefaultContextManager.Context(app, w, r)
	ctx = acontext.WithVars(ctx, r)
	ctx = acontext.WithLogger(ctx, acontext.GetLogger(ctx))
	arc := &appRequestContext{
		Context: ctx,
	}

	arc.URLBuilder = v1.NewURLBuilderFromRequest(r, false)
	return arc
}

func (app *App) register(routeName string, dispatch dispatchFunc) {
	app.router.GetRoute(routeName).Handler(app.dispatcher(dispatch))
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := acontext.DefaultContextManager.Context(app, w, r)
	defer func() {
		status, ok := ctx.Value("http.response.status").(int)
		if ok && status >= 200 && status <= 399 {
			acontext.GetResponseLogger(ctx).Infof("response completed")
		}
	}()

	var err error
	w, err = acontext.GetResponseWriter(ctx)
	if err != nil {
		acontext.GetLogger(ctx).Warnf("response writer not found in context")
	}

	w.Header().Add("TINKERSNEST-VERSION", acontext.GetVersion(ctx))
	app.router.ServeHTTP(w, r)
}

func preloadClaim(ctx *appRequestContext, r *http.Request) error {
	code := r.Header.Get("TINKERSNEST-CLAIM")
	if code != "" {
		route := mux.CurrentRoute(r)
		routeName := route.GetName()

		rclaim, err := cqrs.DispatchQuery(ctx, &queries.FindClaim{Code: code})
		if err != nil && err != storage.ErrNotFound {
			return err
		} else if err == storage.ErrNotFound {
			return fmt.Errorf("no such claim")
		}

		claim, ok := rclaim.(*v1.Claim)
		if !ok {
			// TODO: api error type
			return fmt.Errorf("claim data is invalid")
		}

		if claim.Redeemed != 0 {
			// TODO: api error type
			return fmt.Errorf("no such claim")
		}

		ctx.Context = context.WithValue(ctx.Context, "claim", claim)
		ctx.Context = acontext.WithLogger(ctx.Context, acontext.GetLoggerWithField(ctx.Context, "claim", code))

		expect := v1.NoResource
		switch routeName {
		case v1.RouteNameBlog:
			expect = v1.PostResource
		case v1.RouteNameUserRegistry:
			expect = v1.UserResource
		default:
			// not something that requires a claim
			acontext.GetLogger(ctx).Warn("ignoring unneeded claim for this request")
		}

		if expect != v1.NoResource && expect != claim.ResourceType {
			// TODO: api error type
			return fmt.Errorf("claim cannot be used for %q resources", expect)
		}
	}

	return nil
}

func apiBase(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{}"

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, emptyJSON)
}

func withTraceLogging(name string, h func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := acontext.DefaultContextManager.Context(nil, w, r)
		acontext.GetLogger(ctx).Debugf("%s begin", name)
		defer acontext.GetLogger(ctx).Debugf("%s end", name)
		h(w, r)
	})
}
