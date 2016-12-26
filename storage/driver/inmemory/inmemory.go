package inmemory

import (
	"github.com/danielkrainas/gobag/decouple/drivers"

	"github.com/danielkrainas/tinkersnest/storage"
	"github.com/danielkrainas/tinkersnest/storage/driver/factory"
)

type driverFactory struct{}

func (f *driverFactory) Create(parameters map[string]interface{}) (drivers.DriverBase, error) {
	return &driver{
		stores: make(map[string]interface{}, 0),
	}, nil
}

func init() {
	factory.Register("inmemory", &driverFactory{})
}

type driver struct {
	stores map[string]interface{}
}

func (d *driver) Users() storage.UserStore {
	store, ok := d.stores["user"].(storage.UserStore)
	if !ok {
		store = &userStore{}
		d.stores["user"] = store
	}

	return store
}

func (d *driver) Claims() storage.ClaimStore {
	store, ok := d.stores["claim"].(storage.ClaimStore)
	if !ok {
		store = &claimStore{}
		d.stores["claim"] = store
	}

	return store
}

func (d *driver) Posts() storage.PostStore {
	store, ok := d.stores["post"].(storage.PostStore)
	if !ok {
		store = &postStore{}
		d.stores["post"] = store
	}

	return store
}
