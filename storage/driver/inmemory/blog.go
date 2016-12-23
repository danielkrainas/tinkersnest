package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/danielkrainas/gobag/decouple/cqrs"
	"github.com/danielkrainas/gobag/util/slugify"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/commands"
	"github.com/danielkrainas/tinkersnest/queries"
	"github.com/danielkrainas/tinkersnest/storage"
)

var blog *postStore

func init() {
	blog = &postStore{
		posts: make([]*v1.Post, 0),
	}

	registerCommand(&commands.StorePost{}, blog)
	registerCommand(&commands.DeletePost{}, blog)
	registerQuery(&queries.SearchPosts{}, blog)
	registerQuery(&queries.FindPost{}, blog)
}

type postStore struct {
	m     sync.Mutex
	posts []*v1.Post
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
	s.m.Lock()
	defer s.m.Unlock()
	return s.posts[:], nil
}

func (s *postStore) FindPost(ctx context.Context, q *queries.FindPost) (interface{}, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, p := range s.posts {
		if p.Name == q.Name {
			return p, nil
		}
	}

	return nil, nil
}

func (s *postStore) DeletePost(ctx context.Context, c *commands.DeletePost) error {
	s.m.Lock()
	defer s.m.Unlock()
	for i, p := range s.posts {
		if p.Name == c.Name {
			s.posts = append(s.posts[:i], s.posts[i+1:]...)
			return nil
		}
	}

	return storage.ErrNotFound
}

func (s *postStore) StorePost(ctx context.Context, c *commands.StorePost) error {
	p := c.Post
	if c.New {
		p.Created = time.Now().Unix()
	}

	if p.Name == "" {
		p.Name = slugify.Marshal(p.Title)
	}

	s.m.Lock()
	defer s.m.Unlock()

	found := false
	if !c.New {
		for i, p2 := range s.posts {
			if p2.Name == p.Name {
				s.posts[i] = p
				found = true
				break
			}
		}
	}

	if !found {
		s.posts = append(s.posts, p)
	}

	return nil
}
