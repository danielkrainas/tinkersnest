package mongodb

import "gopkg.in/mgo.v2/bson"

func nameQuery(name string) bson.M {
	return bson.M{"name": name}
}
