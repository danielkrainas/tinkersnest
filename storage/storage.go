package storage

import (
	"context"
	"errors"

	"github.com/danielkrainas/gobag/decouple/cqrs"
	"github.com/danielkrainas/gobag/decouple/drivers"
)

var ErrNotFound = errors.New("not found")

type Driver interface {
	drivers.DriverBase

	Users() UserStore
	Claims() ClaimStore
	Posts() PostStore
}

type UserStore interface {
	Find(name string) (*v1.User, error)
	Count(f *UserFilters) (int, error)
}

type ClaimStore interface {
	Find(code string) (*v1.Claim, error)
}

type PostStore interface {
	Delete(p *v1.Post) error
	Store(p *v1.Post, isNew bool) error
	Find(name string) (*v1.Post, error)
	FindMany(f *PostFilters) ([]*v1.Post, error)
}

type PostFilters struct{}

type UserFilters struct{}
