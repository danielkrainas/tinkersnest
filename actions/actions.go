package actions

import (
	"github.com/danielkrainas/gobag/decouple/cqrs"

	"github.com/danielkrainas/tinkernest/api/v1"
	"github.com/danielkrainas/tinkernest/configuration"
	"github.com/danielkrainas/tinkernest/storage"
)

func FindClaim(ctx context.Context, q *queries.FindClaim, claims storage.ClaimStore) (*v1.Claim, error) {

}

func RedeemClaim(ctx context.Context, c *commands.RedeemClaim, claims storage.ClaimStore) error {

}

func CreateClaim(ctx context.Context, c *commands.CreateClaim, claims storage.ClaimStore) error {

}

func DeleteUser(ctx context.Context, c *commands.DeleteUser, users storage.UserStore) error {

}

func StoreUser(ctx context.Context, c *commands.StoreUser, users storage.UserStore) error {

}

func FindUser(ctx context.Context, q *queries.FindUser, users storage.UserStore) (*v1.User, error) {

}

func CountUsers(ctx context.Context, q *queries.CountUsers, users storage.UserStore) (int, error) {

}

func StorePost(ctx context.Context, c *commands.StorePost, posts storage.PostStore) error {

}

func DeletePost(ctx context.Context, c *commands.DeletePost, posts storage.PostStore) error {

}

func SearchPosts(ctx context.Context, q *queries.SearchPosts, posts storage.PostStore) ([]*v1.Post, error) {

}

func FindPost(ctx context.Context, q *queries.FindPost, posts storage.PostStore) (*v1.Post, error) {

}
