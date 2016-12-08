package commands

import (
	"github.com/danielkrainas/tinkersnest/api/v1"
)

type StorePost struct {
	New  bool
	Post *v1.Post
}

type DeletePost struct {
	Name string
}

type CreateClaim struct {
	Code         string
	ResourceType v1.ResourceType
}

type RedeemClaim struct {
	Code         string
	ResourceType v1.ResourceType
}

type StoreUser struct {
	New  bool
	User *v1.User
}
