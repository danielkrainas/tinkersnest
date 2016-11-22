package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/danielkrainas/tinkersnest/api/errcode"
	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/context"
	"github.com/danielkrainas/tinkersnest/cqrs"
	"github.com/danielkrainas/tinkersnest/cqrs/commands"
	"github.com/danielkrainas/tinkersnest/cqrs/queries"
)

func userRegistryDispatcher(ctx context.Context, r *http.Request) http.Handler {
	h := &userHandler{
		Context: ctx,
	}

	return handlers.MethodHandler{
		"GET":  withTraceLogging("GetAllUsers", h.GetAllUsers),
		"POST": withTraceLogging("CreateUser", h.CreateUser),
	}
}

func userByNameDispatcher(ctx context.Context, r *http.Request) http.Handler {
	h := &userHandler{
		Context: ctx,
	}

	return handlers.MethodHandler{
		"GET": withTraceLogging("GetUser", h.GetUser),
	}
}

type userHandler struct {
	context.Context
}

func (ctx *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userName := acontext.GetStringValue(ctx, "vars.user_name")
	post, err := cqrs.DispatchQuery(ctx, &queries.FindUser{
		Name: userName,
	})

	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if post == nil {
		ctx.Context = acontext.AppendError(ctx, v1.ErrorCodeResourceUnknown)
		return
	}

	if err := v1.ServeJSON(w, post); err != nil {
		acontext.GetLogger(ctx).Errorf("error sending post json: %v", err)
	}
}

func (ctx *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	u := &v1.User{}
	if err = json.Unmarshal(body, u); err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if err := cqrs.DispatchCommand(ctx, &commands.StoreUser{New: true, User: u}); err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if claim, ok := ctx.Value("claim").(*v1.Claim); ok {
		err := cqrs.DispatchCommand(ctx, &commands.RedeemClaim{Code: claim.Code})
		if err != nil {
			acontext.GetLogger(ctx).Error(err)
			ctx.Context = acontext.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
			return
		}
	}

	acontext.GetLoggerWithField(ctx, "user.name", u.Name).Infof("user %q created", u.Name)
	if err := v1.ServeJSON(w, u); err != nil {
		acontext.GetLogger(ctx).Errorf("error sending user json: %v", err)
	}
}

func (ctx *userHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cqrs.DispatchQuery(ctx, &queries.SearchUsers{})
	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if err := v1.ServeJSON(w, users); err != nil {
		acontext.GetLogger(ctx).Errorf("error sending users json: %v", err)
	}
}
