package mongodb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/storage"
)

const usersCollection = "users"

func newUserStore(driver *driver) *userStore {
	store := &userStore{driver.db}
	return store
}

type userStore struct {
	db *mgo.Database
}

var _ storage.UserStore = &userStore{}

func (s *userStore) Delete(name string) error {
	return s.db.C(usersCollection).Remove(nameQuery(name))
}

func (s *userStore) Store(u *v1.User, isNew bool) error {
	users := s.db.C(usersCollection)
	_, err := users.Upsert(nameQuery(u.Name), bson.M{"$set": u})
	return err
}

func (s *userStore) Find(name string) (*v1.User, error) {
	u := &v1.User{}
	iter := s.db.C(usersCollection).Find(nameQuery(name)).Iter()
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

func (s *userStore) FindMany(f *storage.UserFilters) ([]*v1.User, error) {
	users := make([]*v1.User, 0)
	iter := s.db.C(usersCollection).Find(bson.M{}).Iter()
	user := v1.User{}
	for iter.Next(&user) {
		u := user
		users = append(users, &u)
	}

	if iter.Err() != nil {
		return nil, iter.Err()
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return users, nil

}

func (s *userStore) Count(f *storage.UserFilters) (int, error) {
	return s.db.C(usersCollection).Count()
}
