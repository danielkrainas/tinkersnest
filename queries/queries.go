package queries

import (
	"github.com/danielkrainas/tinkersnest/api/v1"
)

type SearchPosts struct {
	Author *v1.Author
}

type FindPost struct {
	Name string
}

type CountUsers struct{}

type FindClaim struct {
	Code string
}

type FindUser struct {
	Name string
}

type SearchUsers struct{}
