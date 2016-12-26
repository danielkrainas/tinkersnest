package mongodb

import (
	"errors"

	"github.com/danielkrainas/gobag/decouple/drivers"
	"gopkg.in/mgo.v2"

	"github.com/danielkrainas/tinkersnest/storage"
	"github.com/danielkrainas/tinkersnest/storage/driver/factory"
)

const (
	postsCollection  = "posts"
	claimsCollection = "claims"
	usersCollection  = "users"
)

type driverFactory struct{}

func (f *driverFactory) Create(parameters map[string]interface{}) (drivers.DriverBase, error) {
	url, ok := parameters["url"].(string)
	if !ok || url == "" {
		return nil, errors.New("url parameter invalid or missing")
	}

	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	d := &driver{
		session: session,
		db:      session.DB(""),
	}

	if err := d.Init(); err != nil {
		return nil, err
	}

	return d, nil
}

func init() {
	factory.Register("mongodb", &driverFactory{})
}

type driver struct {
	session *mgo.Session
	db      *mgo.Database

	users  *userStore
	posts  *postStore
	claims *claimStore
}

var _ storage.Driver = &driver{}

func (d *driver) Init() error {
	d.users = &userStore{d.db}
	d.posts = &postStore{d.db}
	d.claims = &claimStore{d.db}

	nameIndex := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	}

	d.db.C(postsCollection).EnsureIndex(nameIndex)
	d.db.C(usersCollection).EnsureIndex(nameIndex)
	d.db.C(claimsCollection).EnsureIndex(mgo.Index{
		Key:        []string{"code"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	})

	return nil
}

func (d *driver) Users() storage.UserStore {
	return d.users
}

func (d *driver) Posts() storage.PostStore {
	return d.posts
}

func (d *driver) Claims() storage.ClaimStore {
	return d.claims
}
