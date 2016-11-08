package inmemory

import (
	"github.com/danielkrainas/tinkersnest/cqrs"
	"github.com/danielkrainas/tinkersnest/storage"
	"github.com/danielkrainas/tinkersnest/storage/factory"
)

var (
	queryRouter = &cqrs.QueryRouter{}

	commandRouter = &cqrs.CommandRouter{}
)

func registerQuery(q cqrs.Query, exec cqrs.QueryExecutor) {
	queryRouter.Register(q, exec)
}

func registerCommand(c cqrs.Command, handler cqrs.CommandHandler) {
	commandRouter.Register(c, handler)
}

type driverFactory struct{}

func (f *driverFactory) Create(parameters map[string]interface{}) (storage.Driver, error) {
	return &driver{}, nil
}

func init() {
	factory.Register("inmemory", &driverFactory{})
}

type driver struct{}

func (d *driver) Init() error {
	return nil
}

func (d *driver) Command() cqrs.CommandHandler {
	return commandRouter
}

func (d *driver) Query() cqrs.QueryExecutor {
	return queryRouter
}
