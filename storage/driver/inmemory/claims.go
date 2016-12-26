package inmemory

import (
	"sync"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/storage"
)

type claimStore struct {
	m      sync.Mutex
	claims []*v1.Claim
}

func (s *claimStore) Store(c *v1.Claim, isNew bool) error {
	s.m.Lock()
	defer s.m.Unlock()
	s.claims = append(s.claims, c)
	return nil
}

func (s *claimStore) Find(code string) (*v1.Claim, error) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, c := range s.claims {
		if c.Code == code {
			return c, nil
		}
	}

	return nil, storage.ErrNotFound
}
