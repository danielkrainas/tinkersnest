package inmemory

import (
	"sync"
	"time"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/context"
	"github.com/danielkrainas/tinkersnest/cqrs"
	"github.com/danielkrainas/tinkersnest/cqrs/commands"
	"github.com/danielkrainas/tinkersnest/cqrs/queries"
	"github.com/danielkrainas/tinkersnest/util/slugify"
	"github.com/danielkrainas/tinkersnest/util/uuid"
)

var blog *postStore

func init() {
	blog = &postStore{
		posts: make([]*v1.Post, 0),
	}

	registerCommand(&commands.StorePost{}, blog)
	registerQuery(&queries.SearchPosts{}, blog)
}

type postStore struct {
	m     sync.Mutex
	posts []*v1.Post
}

func (s *postStore) Execute(ctx context.Context, q cqrs.Query) (interface{}, error) {
	switch q := q.(type) {
	case *queries.SearchPosts:
		return s.SearchPosts(ctx, q)
	}

	return nil, cqrs.ErrNoExecutor
}

func (s *postStore) Handle(ctx context.Context, c cqrs.Command) error {
	switch c := c.(type) {
	case *commands.StorePost:
		return s.StorePost(ctx, c)
	}

	return cqrs.ErrNoHandler
}

func (s *postStore) SearchPosts(ctx context.Context, q *queries.SearchPosts) (interface{}, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.posts[:], nil
}

func (s *postStore) StorePost(ctx context.Context, c *commands.StorePost) error {
	p := c.Post
	if c.New || p.ID == "" {
		p.ID = uuid.Generate()
		p.Created = time.Now().Unix()
	}

	if p.Slug == "" {
		p.Slug = slugify.Marshal(p.Title)
	}

	s.m.Lock()
	defer s.m.Unlock()

	found := false
	if !c.New {
		for i, p2 := range s.posts {
			if p2.ID == p.ID {
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
