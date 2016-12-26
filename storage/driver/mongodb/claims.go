package mongodb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/storage"
)

const claimsCollection = "claims"

func newClaimStore(driver *driver) *claimStore {
	store := &claimStore{driver.db}
	return store
}

type claimStore struct {
	db *mgo.Database
}

func (s *claimStore) Store(c *v1.Claim, isNew bool) error {
	claims := s.db.C(claimsCollection)
	_, err := claims.Upsert(bson.M{"code": c.Code}, bson.M{"$set": c})
	return err
}

func (s *claimStore) Find(code string) (*v1.Claim, error) {
	c := &v1.Claim{}
	iter := s.db.C(claimsCollection).Find(bson.M{"code": code}).Iter()
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
