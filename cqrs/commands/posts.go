package commands

import (
	"github.com/danielkrainas/tinkersnest/api/v1"
)

type StorePost struct {
	New  bool
	Post *v1.Post
}
