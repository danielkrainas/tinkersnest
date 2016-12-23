package factory

import (
	"github.com/danielkrainas/gobag/decouple/drivers"

	"github.com/danielkrainas/tinkersnest/storage/driver"
)

const assetType = "Storage"

var registry = &drivers.Registry{
	AssetType: "Storage",
}

func Register(name string, factory drivers.Factory) {
	registry.Register(name, factory)
}

func Create(name string, parameters map[string]interface{}) (driver.Driver, error) {
	d, err := registry.Create(name, parameters)
	return d.(driver.Driver), err
}
