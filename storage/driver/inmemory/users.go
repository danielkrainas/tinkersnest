package inmemory

import (
	"sync"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/storage"
)

type userStore struct {
	m     sync.Mutex
	users []*v1.User
}

func (s *userStore) Store(u *v1.User, isNew bool) error {
	s.m.Lock()
	defer s.m.Unlock()

	found := false
	if isNew {
		for i, u2 := range s.users {
			if u2.Name == u.Name {
				s.users[i] = u
				found = true
				break
			}
		}
	}

	if !found {
		s.users = append(s.users, u)
	}

	return nil
}

func (s *userStore) Delete(name string) error {
	s.m.Lock()
	defer s.m.Unlock()
	for i, u := range s.users {
		if u.Name == name {
			s.users = append(s.users[:i], s.users[i+1:]...)
			return nil
		}
	}

	return storage.ErrNotFound
}

func (s *userStore) Find(name string) (*v1.User, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, u := range s.users {
		if u.Name == name {
			return u, nil
		}
	}

	return nil, nil
}

func (s *userStore) Count(f *storage.UserFilters) (int, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return len(s.users), nil
}

func (s *userStore) FindMany(f *storage.UserFilters) ([]*v1.User, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.users[:], nil
}
