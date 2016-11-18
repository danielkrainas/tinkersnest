package storage

import (
	"github.com/danielkrainas/tinkersnest/configuration"
	"github.com/danielkrainas/tinkersnest/context"
	"github.com/danielkrainas/tinkersnest/storage/driver"
	"github.com/danielkrainas/tinkersnest/storage/driver/factory"
)

func FromConfig(config *configuration.Config) (driver.Driver, error) {
	params := config.Storage.Parameters()
	if params == nil {
		params = make(configuration.Parameters)
	}

	d, err := factory.Create(config.Storage.Type(), params)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func LogSummary(ctx context.Context, config *configuration.Config) {
	context.GetLogger(ctx).Infof("using %q storage driver", config.Storage.Type())
}
