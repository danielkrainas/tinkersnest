package handlers

import (
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
		"GET":  http.HandlerFunc(h.GetAllPosts),
		"POST": http.HandlerFunc(h.CreatePost),
	}
}

type blogHandler struct {
	context.Context
}

func (ctx *blogHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	context.GetLogger(ctx).Debug("CreatePost begin")
	defer context.GetLogger(ctx).Debug("CreatePost end")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	p := &v1.Post{}
	if err = json.Unmarshal(body, p); err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if err := cqrs.DispatchCommand(ctx, &commands.StorePost{New: true, Post: p}); err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	context.GetLoggerWithField(ctx, "post.id", p.ID).Infof("blog post %q created")
	if err := v1.ServeJSON(w, p); err != nil {
		context.GetLogger(ctx).Errorf("error sending blog post json: %v", err)
	}
}

func (ctx *blogHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	context.GetLogger(ctx).Debug("GetAllPosts begin")
	defer context.GetLogger(ctx).Debug("GetAllPosts end")

	posts, err := cqrs.DispatchQuery(ctx, &queries.SearchPosts{})
	if err != nil {
		context.GetLogger(ctx).Error(err)
		ctx.Context = context.AppendError(ctx.Context, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	if err := v1.ServeJSON(w, posts); err != nil {
		context.GetLogger(ctx).Errorf("error sending posts json: %v", err)
	}
}
