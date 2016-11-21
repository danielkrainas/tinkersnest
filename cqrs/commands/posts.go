package commands

import (
	"github.com/danielkrainas/tinkersnest/api/v1"
)

type StorePost struct {
	New  bool
	Post *v1.Post
}

type CreateClaim struct {
	Code         string
	ResourceType v1.ResourceType
}

type RedeemClaim struct {
	Code string
}

type StoreUser struct {
	User *v1.User
}
