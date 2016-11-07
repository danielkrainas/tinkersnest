package cqrs

import (
	"errors"
	"fmt"

	"github.com/danielkrainas/canaria-api/context"
)

var ErrNoHandler = errors.New("no command handler")

type Command interface{}

type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

type CommandDispatcher struct {
	Handlers []CommandHandler
}

func (d *CommandDispatcher) Dispatch(ctx context.Context, cmd Command) error {
	for _, h := range d.Handlers {
		if err := h.Handle(ctx, cmd); err != nil && err != ErrNoHandler {
			return err
		} else if err == nil {
			return nil
		}
	}

	return ErrNoHandler
}

type CommandRouter map[string]CommandHandler

func getCommandKey(c Command) string {
	return fmt.Sprintf("%T", c)
}

func (r CommandRouter) Register(c Command, exec CommandHandler) {
	r[getCommandKey(c)] = exec
}

func (r CommandRouter) Handle(ctx context.Context, c Command) error {
	h, ok := r[getCommandKey(c)]
	if !ok {
		return ErrNoHandler
	}

	return h.Handle(ctx, c)
}

func WithCommandDispatch(ctx context.Context, d *CommandDispatcher) context.Context {
	return context.WithValue(ctx, "cmd.dispatcher", d)
}

func DispatchCommand(ctx context.Context, c Command) error {
	d, ok := ctx.Value("cmd.dispatcher").(*CommandDispatcher)
	if !ok || d == nil {
		return fmt.Errorf("no valid command dispatchers found in context")
	}

	return d.Dispatch(ctx, c)
}
