package handlers

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/danielkrainas/tinkersnest/api/errcode"
	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/auth"
	"github.com/danielkrainas/tinkersnest/context"
	"github.com/danielkrainas/tinkersnest/cqrs"
	"github.com/danielkrainas/tinkersnest/cqrs/queries"
)

func authDispatcher(ctx context.Context, r *http.Request) http.Handler {
	h := &authHandler{
		Context: ctx,
	}

	return handlers.MethodHandler{
		"POST": withTraceLogging("Authorize", h.Auth),
	}
}

type authHandler struct {
	context.Context
}

func (ctx *authHandler) Auth(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	creds := &v1.User{}
	if err = json.Unmarshal(body, creds); err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	userData, err := cqrs.DispatchQuery(ctx, &queries.FindUser{Name: creds.Name})
	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	user, ok := userData.(*v1.User)
	if !ok {
		acontext.GetLogger(ctx).Error("couldn't cast user data")
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown)
		return
	}

	creds.HashedPassword = auth.HashPassword(creds.Password, user.Salt)
	if creds.HashedPassword != creds.Password {
		acontext.GetLogger(ctx).Error("invalid username or password")
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail("invalid username or password"))
		return
	}

	token, err := auth.BearerToken(user)
	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if _, err = io.WriteString(w, token); err != nil {
		acontext.GetLogger(ctx).Errorf("error sending auth token: %v", err)
	}
}
