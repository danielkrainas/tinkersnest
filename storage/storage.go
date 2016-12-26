package storage

import (
	"errors"

	"github.com/danielkrainas/gobag/decouple/drivers"

	"github.com/danielkrainas/tinkersnest/api/v1"
)

var ErrNotFound = errors.New("not found")

type Driver interface {
	drivers.DriverBase

	Users() UserStore
	Claims() ClaimStore
	Posts() PostStore
}

type UserStore interface {
	Delete(name string) error
	Store(u *v1.User, isNew bool) error
	Find(name string) (*v1.User, error)
	FindMany(f *UserFilters) ([]*v1.User, error)
	Count(f *UserFilters) (int, error)
}

type ClaimStore interface {
	Find(code string) (*v1.Claim, error)
	Store(c *v1.Claim, isNew bool) error
}

type PostStore interface {
	Delete(name string) error
	Store(p *v1.Post, isNew bool) error
	Find(name string) (*v1.Post, error)
	FindMany(f *PostFilters) ([]*v1.Post, error)
}

type PostFilters struct{}

type UserFilters struct{}
