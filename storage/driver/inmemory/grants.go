package inmemory

import (
	"context"
	"crypto/md5"
	"fmt"
	"sync"
	"time"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/cqrs"
	"github.com/danielkrainas/tinkersnest/cqrs/commands"
	"github.com/danielkrainas/tinkersnest/cqrs/queries"
	"github.com/danielkrainas/tinkersnest/util/uuid"
)

var grants *grantStore

func init() {
	grants = &grantStore{
		grants: make([]*v1.Grant, 0),
	}

	registerCommand(&commands.StoreGrant{}, grants)
	registerQuery(&queries.SearchGrants{}, grants)
	registerQuery(&queries.FindGrant{}, grants)
}

type grantStore struct {
	m      sync.Mutex
	grants []*v1.Grant
}

func (s *grantStore) Execute(ctx context.Context, q cqrs.Query) (interface{}, error) {
	switch q := q.(type) {
	case *queries.SearchGrants:
		return s.SearchGrants(ctx, q)

	case *queries.FindGrant:
		return s.FindGrant(ctx, q)
	}

	return nil, cqrs.ErrNoExecutor
}

func (s *grantStore) Handle(ctx context.Context, c cqrs.Command) error {
	switch c := c.(type) {
	case *commands.StoreGrant:
		return s.StoreGrant(ctx, c)
	}

	return cqrs.ErrNoHandler
}

func (s *grantStore) SearchGrants(ctx context.Context, q *queries.SearchGrants) (interface{}, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.grants[:], nil
}

func (s *grantStore) FindGrant(ctx context.Context, q *queries.FindGrant) (interface{}, error) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, p := range s.grants {
		if p.Code == q.Code {
			return p, nil
		}
	}

	return nil, nil
}

func (s *grantStore) StoreGrant(ctx context.Context, c *commands.StoreGrant) error {
	g := c.Grant
	if c.New {
		g.Created = time.Now().Unix()
		g.Code = newGrantCode(g)
	}

	s.m.Lock()
	defer s.m.Unlock()

	found := false
	if !c.New {
		for i, g2 := range s.grants {
			if g2.Code == g.Code {
				s.grants[i] = g
				found = true
				break
			}
		}
	}

	if !found {
		s.grants = append(s.grants, g)
	}

	return nil
}

func newGrantCode(g *v1.Grant) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s.%d.%s", uuid.Generate(), time.Now().Unix(), g.ResourceType))))
}
