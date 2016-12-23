package storage

import (
	"context"
	"errors"

	cfg "github.com/danielkrainas/gobag/configuration"
	"github.com/danielkrainas/gobag/context"

	"github.com/danielkrainas/tinkersnest/configuration"
	"github.com/danielkrainas/tinkersnest/storage/driver"
	"github.com/danielkrainas/tinkersnest/storage/driver/factory"
)

var (
	ErrNotFound = errors.New("not found")
)

func FromConfig(config *configuration.Config) (driver.Driver, error) {
	params := config.Storage.Parameters()
	if params == nil {
		params = make(cfg.Parameters)
	}

	d, err := factory.Create(config.Storage.Type(), params)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func LogSummary(ctx context.Context, config *configuration.Config) {
	acontext.GetLogger(ctx).Infof("using %q storage driver", config.Storage.Type())
}
