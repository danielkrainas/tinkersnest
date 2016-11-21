package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/cqrs"
	"github.com/danielkrainas/tinkersnest/cqrs/commands"
	"github.com/danielkrainas/tinkersnest/cqrs/queries"
	"github.com/danielkrainas/tinkersnest/storage"
)

var claims *claimStore

func init() {
	claims = &claimStore{
		claims: make([]*v1.Claim, 0),
	}

	registerCommand(&commands.CreateClaim{}, claims)
	registerCommand(&commands.RedeemClaim{}, claims)
	registerQuery(&queries.FindClaim{}, claims)
}

type claimStore struct {
	m      sync.Mutex
	claims []*v1.Claim
}

func (s *claimStore) Execute(ctx context.Context, q cqrs.Query) (interface{}, error) {
	switch q := q.(type) {
	case *queries.FindClaim:
		return s.FindClaim(ctx, q)
	}

	return nil, cqrs.ErrNoExecutor
}

func (s *claimStore) Handle(ctx context.Context, c cqrs.Command) error {
	switch c := c.(type) {
	case *commands.CreateClaim:
		return s.CreateClaim(ctx, c)
	case *commands.RedeemClaim:
		return s.RedeemClaim(ctx, c)
	}

	return cqrs.ErrNoHandler
}

func (s *claimStore) CreateClaim(ctx context.Context, c *commands.CreateClaim) error {
	claim := &v1.Claim{
		Code:         c.Code,
		ResourceType: c.ResourceType,
		Created:      time.Now().Unix(),
		Redeemed:     0,
	}

	s.m.Lock()
	defer s.m.Unlock()
	s.claims = append(s.claims, claim)
	return nil
}

func (s *claimStore) RedeemClaim(ctx context.Context, c *commands.RedeemClaim) error {
	s.m.Lock()
	defer s.m.Unlock()
	for _, claim := range s.claims {
		if claim.Code == c.Code && claim.Redeemed < 1 {
			claim.Redeemed = time.Now().Unix()
			return nil
		}
	}

	return storage.ErrNotFound
}

func (s *claimStore) FindClaim(ctx context.Context, q *queries.FindClaim) (interface{}, error) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, c := range s.claims {
		if c.Code == q.Code {
			return c, nil
		}
	}

	return nil, nil
}
