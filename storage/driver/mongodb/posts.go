package mongodb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/storage"
)

const postsCollection = "posts"

func newPostStore(driver *driver) *postStore {
	store := &postStore{driver.db}
	return store
}

type postStore struct {
	db *mgo.Database
}

var _ storage.PostStore = &postStore{}

func (s *postStore) Delete(name string) error {
	return s.db.C(postsCollection).Remove(nameQuery(name))
}

func (s *postStore) Store(p *v1.Post, isNew bool) error {
	posts := s.db.C(postsCollection)
	_, err := posts.Upsert(nameQuery(p.Name), bson.M{"$set": p})
	return err
}

func (s *postStore) Find(name string) (*v1.Post, error) {
	p := &v1.Post{}
	iter := s.db.C(postsCollection).Find(nameQuery(name)).Iter()
	if !iter.Next(p) {
		return nil, storage.ErrNotFound
	}

	if iter.Err() != nil {
		return nil, iter.Err()
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *postStore) FindMany(f *storage.PostFilters) ([]*v1.Post, error) {
	posts := make([]*v1.Post, 0)
	iter := s.db.C(postsCollection).Find(bson.M{}).Iter()
	post := v1.Post{}
	for iter.Next(&post) {
		p := post
		posts = append(posts, &p)
	}

	if iter.Err() != nil {
		return nil, iter.Err()
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return posts, nil
}
