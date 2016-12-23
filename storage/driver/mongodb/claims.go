package mongodb

import (
	"context"
	"time"

	"github.com/danielkrainas/gobag/decouple/cqrs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/commands"
	"github.com/danielkrainas/tinkersnest/queries"
	"github.com/danielkrainas/tinkersnest/storage"
)

const claimsCollection = "claims"

func newClaimStore(driver *driver) *claimStore {
	store := &claimStore{driver.db}
	driver.registerCommand(&commands.CreateClaim{}, store)
	driver.registerCommand(&commands.RedeemClaim{}, store)
	driver.registerQuery(&queries.FindClaim{}, store)
	return store
}

type claimStore struct {
	db *mgo.Database
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

	claims := s.db.C(claimsCollection)
	return claims.Insert(claim)
}

func (s *claimStore) RedeemClaim(ctx context.Context, c *commands.RedeemClaim) error {
	return s.db.C(claimsCollection).Update(bson.M{
		"code":     c.Code,
		"redeemed": 0,
	}, bson.M{
		"$set": bson.M{
			"redeemed": time.Now().Unix(),
		},
	})
}

func (s *claimStore) FindClaim(ctx context.Context, q *queries.FindClaim) (interface{}, error) {
	c := &v1.Claim{}
	iter := s.db.C(claimsCollection).Find(bson.M{"code": q.Code}).Iter()
	if !iter.Next(c) {
		return nil, storage.ErrNotFound
	}

	if iter.Err() != nil {
		return nil, iter.Err()
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return c, nil
}
