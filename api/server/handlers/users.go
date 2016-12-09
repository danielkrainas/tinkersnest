package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/danielkrainas/tinkersnest/api/errcode"
	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/auth"
	"github.com/danielkrainas/tinkersnest/context"
	"github.com/danielkrainas/tinkersnest/cqrs"
	"github.com/danielkrainas/tinkersnest/cqrs/commands"
	"github.com/danielkrainas/tinkersnest/cqrs/queries"
	"github.com/danielkrainas/tinkersnest/storage"
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
		"DELETE": withTraceLogging("DeleteUser", h.DeleteUser),
		"PUT": withTraceLogging("UpdateUser", h.UpdateUser),
	}
}

type userHandler struct {
	context.Context
}

func (ctx *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userName := acontext.GetStringValue(ctx, "vars.user_name")
	userRaw, err := cqrs.DispatchQuery(ctx, &queries.FindUser{
		Name: userName,
	})

	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if userRaw == nil {
		ctx.Context = acontext.AppendError(ctx, v1.ErrorCodeResourceUnknown)
		return
	}

	user := userRaw.(*v1.User)
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

	if u.Password != "" {
		user.HashedPassword = auth.HashPassword(u.Password, user.Salt)
	}

	if u.Email != "" {
		user.Email = u.Email
	}

	if u.FullName != "" {
		user.FullName = u.FullName
	}

	if err := cqrs.DispatchCommand(ctx, &commands.StoreUser{New: false, User: user}); err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if err := v1.ServeJSON(w, user); err != nil {
		acontext.GetLogger(ctx).Errorf("error sending user json: %v", err)
	}
}

func (ctx *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userName := acontext.GetStringValue(ctx, "vars.user_name")
	err := cqrs.DispatchCommand(ctx, &commands.DeleteUser{userName})
	if err != nil {
		if err == storage.ErrNotFound {
			acontext.GetLogger(ctx).Error("user not found")
			ctx.Context = acontext.AppendError(ctx, v1.ErrorCodeResourceUnknown)
			return
		} else {
			acontext.GetLogger(ctx).Error(err)
			ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
			return	
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ctx *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userName := acontext.GetStringValue(ctx, "vars.user_name")
	user, err := cqrs.DispatchQuery(ctx, &queries.FindUser{
		Name: userName,
	})

	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if user == nil {
		ctx.Context = acontext.AppendError(ctx, v1.ErrorCodeResourceUnknown)
		return
	}

	if err := v1.ServeJSON(w, user); err != nil {
		acontext.GetLogger(ctx).Errorf("error sending user json: %v", err)
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

	if u.Salt, err = auth.GenerateSalt(); err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	u.HashedPassword = auth.HashPassword(u.Password, u.Salt)
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
