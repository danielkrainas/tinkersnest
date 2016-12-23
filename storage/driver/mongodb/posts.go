package mongodb

import (
	"context"
	"time"

	"github.com/danielkrainas/gobag/decouple/cqrs"
	"github.com/danielkrainas/gobag/util/slugify"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/commands"
	"github.com/danielkrainas/tinkersnest/queries"
	"github.com/danielkrainas/tinkersnest/storage"
)

const postsCollection = "posts"

func newPostStore(driver *driver) *postStore {
	store := &postStore{driver.db}
	driver.registerCommand(&commands.StorePost{}, store)
	driver.registerCommand(&commands.DeletePost{}, store)
	driver.registerQuery(&queries.SearchPosts{}, store)
	driver.registerQuery(&queries.FindPost{}, store)
	return store
}

type postStore struct {
	db *mgo.Database
}

func (s *postStore) Execute(ctx context.Context, q cqrs.Query) (interface{}, error) {
	switch q := q.(type) {
	case *queries.SearchPosts:
		return s.SearchPosts(ctx, q)

	case *queries.FindPost:
		return s.FindPost(ctx, q)
	}

	return nil, cqrs.ErrNoExecutor
}

func (s *postStore) Handle(ctx context.Context, c cqrs.Command) error {
	switch c := c.(type) {
	case *commands.StorePost:
		return s.StorePost(ctx, c)
	case *commands.DeletePost:
		return s.DeletePost(ctx, c)
	}

	return cqrs.ErrNoHandler
}

func (s *postStore) SearchPosts(ctx context.Context, q *queries.SearchPosts) (interface{}, error) {
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

func (s *postStore) FindPost(ctx context.Context, q *queries.FindPost) (interface{}, error) {
	p := &v1.Post{}
	iter := s.db.C(postsCollection).Find(nameQuery(q.Name)).Iter()
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

func (s *postStore) DeletePost(ctx context.Context, c *commands.DeletePost) error {
	return s.db.C(postsCollection).Remove(nameQuery(c.Name))
}

func (s *postStore) StorePost(ctx context.Context, c *commands.StorePost) error {
	p := c.Post
	if c.New {
		p.Created = time.Now().Unix()
	}

	if p.Name == "" {
		p.Name = slugify.Marshal(p.Title)
	}

	posts := s.db.C(postsCollection)
	_, err := posts.Upsert(nameQuery(p.Name), bson.M{"$set": p})
	return err
}
