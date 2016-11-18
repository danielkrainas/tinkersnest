package factory

import (
	"fmt"

	"github.com/danielkrainas/tinkersnest/storage/driver"
)

const assetType = "Storage"

var factories = make(map[string]DriverFactory)

type DriverFactory interface {
	Create(parameters map[string]interface{}) (driver.Driver, error)
}

func Register(name string, factory DriverFactory) {
	if factory == nil {
		panic(fmt.Sprintf("%s DriverFactory cannot be nil", assetType))
	}

	if _, registered := factories[name]; registered {
		panic(fmt.Sprintf("%s DriverFactory named %s already registered", assetType, name))
	}

	factories[name] = factory
}

func Create(name string, parameters map[string]interface{}) (driver.Driver, error) {
	if factory, ok := factories[name]; ok {
		return factory.Create(parameters)
	}

	return nil, InvalidDriverError{name}
}

type InvalidDriverError struct {
	Name string
}

func (err InvalidDriverError) Error() string {
	return fmt.Sprintf("%s driver not registered: %s", assetType, err.Name)
}
