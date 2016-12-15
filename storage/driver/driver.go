package driver

import (
	"errors"

	"github.com/danielkrainas/tinkersnest/cqrs"
)

var ErrNotSupported = errors.New("the operation is not supported by the driver")

type Driver interface {
	Init() error

	Command() cqrs.CommandHandler
	Query() cqrs.QueryExecutor
}
