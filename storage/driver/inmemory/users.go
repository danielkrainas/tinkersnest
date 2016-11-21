package inmemory

import (
	"context"
	"sync"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/cqrs"
	"github.com/danielkrainas/tinkersnest/cqrs/commands"
	"github.com/danielkrainas/tinkersnest/cqrs/queries"
	"github.com/danielkrainas/tinkersnest/util/slugify"
)

var users *userStore

func init() {
	users = &userStore{
		users: make([]*v1.User, 0),
	}

	registerCommand(&commands.StoreUser{}, users)
	registerQuery(&queries.FindUser{}, users)
}

type userStore struct {
	m     sync.Mutex
	users []*v1.User
}

func (s *userStore) Execute(ctx context.Context, q cqrs.Query) (interface{}, error) {
	switch q := q.(type) {
	case *queries.FindUser:
		return s.FindUser(ctx, q)
	}

	return nil, cqrs.ErrNoExecutor
}

func (s *userStore) Handle(ctx context.Context, c cqrs.Command) error {
	switch c := c.(type) {
	case *commands.StoreUser:
		return s.StoreUser(ctx, c)
	}

	return cqrs.ErrNoHandler
}

func (s *userStore) FindUser(ctx context.Context, q *queries.FindUser) (interface{}, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, p := range s.users {
		if p.Name == q.Name {
			return p, nil
		}
	}

	return nil, nil
}

func (s *userStore) StoreUser(ctx context.Context, c *commands.StoreUser) error {
	u := c.User
	if u.Name == "" {
		u.Name = slugify.Marshal(u.FullName)
	}

	s.m.Lock()
	defer s.m.Unlock()

	found := false
	for i, u2 := range s.users {
		if u2.Name == u.Name {
			s.users[i] = u
			found = true
			break
		}
	}

	if !found {
		s.users = append(s.users, u)
	}

	return nil
}
