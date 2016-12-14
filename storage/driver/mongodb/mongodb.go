package mongodb

import (
	"errors"

	"gopkg.in/mgo.v2"

	"github.com/danielkrainas/tinkersnest/cqrs"
	storage "github.com/danielkrainas/tinkersnest/storage/driver"
	"github.com/danielkrainas/tinkersnest/storage/driver/factory"
)

type driverFactory struct{}

func (f *driverFactory) Create(parameters map[string]interface{}) (storage.Driver, error) {
	url, ok := parameters["url"].(string)
	if !ok || url == "" {
		return nil, errors.New("url parameter invalid or missing")
	}

	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	return &driver{
		session: session,
		db:      session.DB(""),
		qr:      &cqrs.QueryRouter{},
		cr:      &cqrs.CommandRouter{},
	}, nil
}

func init() {
	factory.Register("mongodb", &driverFactory{})
}

type driver struct {
	session *mgo.Session
	db      *mgo.Database

	qr *cqrs.QueryRouter
	cr *cqrs.CommandRouter

	users  *userStore
	posts  *postStore
	claims *claimStore
}

func (d *driver) registerQuery(q cqrs.Query, exec cqrs.QueryExecutor) {
	d.qr.Register(q, exec)
}

func (d *driver) registerCommand(c cqrs.Command, handler cqrs.CommandHandler) {
	d.cr.Register(c, handler)
}

func (d *driver) Init() error {
	d.users = newUserStore(d)
	d.posts = newPostStore(d)
	d.claims = newClaimStore(d)

	nameIndex := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	}

	d.db.C(postsCollection).EnsureIndex(nameIndex)
	return nil
}

func (d *driver) Command() cqrs.CommandHandler {
	return d.cr
}

func (d *driver) Query() cqrs.QueryExecutor {
	return d.qr
}
