package mongodb

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/cqrs"
	"github.com/danielkrainas/tinkersnest/cqrs/commands"
	"github.com/danielkrainas/tinkersnest/cqrs/queries"
	"github.com/danielkrainas/tinkersnest/storage"
	"github.com/danielkrainas/tinkersnest/util/slugify"
)

const usersCollection = "users"

func newUserStore(driver *driver) *userStore {
	store := &userStore{driver.db}
	driver.registerCommand(&commands.DeleteUser{}, store)
	driver.registerCommand(&commands.StoreUser{}, store)
	driver.registerQuery(&queries.FindUser{}, store)
	driver.registerQuery(&queries.CountUsers{}, store)
	return store
}

type userStore struct {
	db *mgo.Database
}

func (s *userStore) Execute(ctx context.Context, q cqrs.Query) (interface{}, error) {
	switch q := q.(type) {
	case *queries.FindUser:
		return s.FindUser(ctx, q)
	case *queries.CountUsers:
		return s.CountUsers(ctx, q)
	}

	return nil, cqrs.ErrNoExecutor
}

func (s *userStore) Handle(ctx context.Context, c cqrs.Command) error {
	switch c := c.(type) {
	case *commands.StoreUser:
		return s.StoreUser(ctx, c)

	case *commands.DeleteUser:
		return s.DeleteUser(ctx, c)
	}

	return cqrs.ErrNoHandler
}

func (s *userStore) CountUsers(ctx context.Context, q *queries.CountUsers) (interface{}, error) {
	return s.db.C(usersCollection).Count()
}

func (s *userStore) FindUser(ctx context.Context, q *queries.FindUser) (interface{}, error) {
	u := &v1.User{}
	iter := s.db.C(usersCollection).Find(nameQuery(q.Name)).Iter()
	if !iter.Next(u) {
		return nil, storage.ErrNotFound
	}

	if iter.Err() != nil {
		return nil, iter.Err()
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userStore) DeleteUser(ctx context.Context, c *commands.DeleteUser) error {
	return s.db.C(usersCollection).Remove(nameQuery(c.Name))
}

func (s *userStore) StoreUser(ctx context.Context, c *commands.StoreUser) error {
	u := c.User
	if u.Name == "" {
		u.Name = slugify.Marshal(u.FullName)
	}

	users := s.db.C(usersCollection)
	_, err := users.Upsert(nameQuery(u.Name), bson.M{"$set": u})
	return err
}
