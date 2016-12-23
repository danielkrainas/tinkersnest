package driver

import (
	"errors"

	"github.com/danielkrainas/gobag/decouple/cqrs"
	"github.com/danielkrainas/gobag/decouple/drivers"
)

var ErrNotSupported = errors.New("the operation is not supported by the driver")

type Driver interface {
	drivers.DriverBase

	Init() error

	Command() cqrs.CommandHandler
	Query() cqrs.QueryExecutor
}
