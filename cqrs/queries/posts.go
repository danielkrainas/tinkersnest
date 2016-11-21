package queries

import ()

type SearchPosts struct{}

type FindPost struct {
	Name string
}

type FindGrant struct {
	Code string
}

type CountUsers struct{}

type FindClaim struct {
	Code string
}

type FindUser struct {
	Name string
}

type SearchUsers struct{}
