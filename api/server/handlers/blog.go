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

func blogListDispatcher(ctx context.Context, r *http.Request) http.Handler {
	h := &blogHandler{
		Context: ctx,
	}

	return handlers.MethodHandler{
		"GET":  withLogging("GetAllPosts", h.GetAllPosts),
		"POST": withLogging("CreatePost", h.CreatePost),
	}
}

type blogHandler struct {
	context.Context
}

func (ctx *blogHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	p := &v1.Post{}
	if err = json.Unmarshal(body, p); err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if err := cqrs.DispatchCommand(ctx, &commands.StorePost{New: true, Post: p}); err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	acontext.GetLoggerWithField(ctx, "post.name", p.Name).Infof("blog post %q created", p.Name)
	if err := v1.ServeJSON(w, p); err != nil {
		acontext.GetLogger(ctx).Errorf("error sending blog post json: %v", err)
	}
}

func (ctx *blogHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := cqrs.DispatchQuery(ctx, &queries.SearchPosts{})
	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		ctx.Context = acontext.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if err := v1.ServeJSON(w, posts); err != nil {
		acontext.GetLogger(ctx).Errorf("error sending posts json: %v", err)
	}
}
