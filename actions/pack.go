package actions

import (
	"github.com/danielkrainas/gobag/decouple/cqrs"

	"github.com/danielkrainas/tinkernest/configuration"
	"github.com/danielkrainas/tinkernest/storage"
	"github.com/danielkrainas/tinkernest/storage/loader"
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
	}

	return p.Query.Execute(ctx, q)
}

func (p *pack) Handle(ctx context.Context, c cqrs.Command) error {
	return p.Command.Handle(ctx, c)
}

func FromConfig(config *configuration.Config) (Pack, error) {
	storageDriver, err := storageloader.FromConfig(config)
	if err != nil {
		return nil, err
	}

	p := &pack{
		store: storageDriver,
	}

	return p
}
