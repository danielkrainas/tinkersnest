package actions

import (
	"context"

	"github.com/danielkrainas/gobag/decouple/cqrs"

	"github.com/danielkrainas/tinkersnest/commands"
	"github.com/danielkrainas/tinkersnest/configuration"
	"github.com/danielkrainas/tinkersnest/queries"
	"github.com/danielkrainas/tinkersnest/storage"
	"github.com/danielkrainas/tinkersnest/storage/loader"
)

type Pack interface {
	cqrs.QueryExecutor
	cqrs.CommandHandler
}

type pack struct {
	store storage.Driver
}

func (p *pack) Execute(ctx context.Context, q cqrs.Query) (interface{}, error) {
	switch q := q.(type) {
	case *queries.FindClaim:
		return FindClaim(ctx, q, p.store.Claims())
	case *queries.FindUser:
		return FindUser(ctx, q, p.store.Users())
	case *queries.CountUsers:
		return CountUsers(ctx, q, p.store.Users())
	case *queries.SearchUsers:
		return SearchUsers(ctx, q, p.store.Users())
	case *queries.SearchPosts:
		return SearchPosts(ctx, q, p.store.Posts())
	case *queries.FindPost:
		return FindPost(ctx, q, p.store.Posts())
	}

	return nil, cqrs.ErrNoExecutor
}

func (p *pack) Handle(ctx context.Context, c cqrs.Command) error {
	switch c := c.(type) {
	case *commands.RedeemClaim:
		return RedeemClaim(ctx, c, p.store.Claims())
	case *commands.CreateClaim:
		return CreateClaim(ctx, c, p.store.Claims())
	case *commands.DeleteUser:
		return DeleteUser(ctx, c, p.store.Users())
	case *commands.StoreUser:
		return StoreUser(ctx, c, p.store.Users())
	case *commands.StorePost:
		return StorePost(ctx, c, p.store.Posts())
	case *commands.DeletePost:
		return DeletePost(ctx, c, p.store.Posts())
	}

	return cqrs.ErrNoHandler
}

func FromConfig(config *configuration.Config) (Pack, error) {
	storageDriver, err := storageloader.FromConfig(config)
	if err != nil {
		return nil, err
	}

	p := &pack{
		store: storageDriver,
	}

	return p, nil
}
