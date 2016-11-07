package storage

import (
	"errors"

	"github.com/danielkrainas/tinkersnest/cqrs"
)

var (
	ErrNotSupported = errors.New("the operation is not supported by the driver")
	ErrNotFound     = errors.New("not found")
)

type Driver interface {
	Command() cqrs.CommandHandler
	Query() cqrs.QueryExecutor
}
