package inmemory

import (
	"sync"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/storage"
)

type postStore struct {
	m     sync.Mutex
	posts []*v1.Post
}

func (s *postStore) FindMany(f *storage.PostFilters) ([]*v1.Post, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.posts[:], nil
}

func (s *postStore) Delete(name string) error {
	s.m.Lock()
	defer s.m.Unlock()
	for i, p := range s.posts {
		if p.Name == name {
			s.posts = append(s.posts[:i], s.posts[i+1:]...)
			return nil
		}
	}

	return storage.ErrNotFound
}

func (s *postStore) Store(p *v1.Post, isNew bool) error {
	s.m.Lock()
	defer s.m.Unlock()

	found := false
	if !isNew {
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

func (s *postStore) Find(name string) (*v1.Post, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, p := range s.posts {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, nil
}
