package actions

import (
	"context"
	"time"

	"github.com/danielkrainas/gobag/util/slugify"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/commands"
	"github.com/danielkrainas/tinkersnest/queries"
	"github.com/danielkrainas/tinkersnest/storage"
)

func FindClaim(ctx context.Context, q *queries.FindClaim, claims storage.ClaimStore) (*v1.Claim, error) {
	return claims.Find(q.Code)
}

func RedeemClaim(ctx context.Context, c *commands.RedeemClaim, claims storage.ClaimStore) error {
	claim, err := claims.Find(c.Code)
	if err != nil {
		return err
	}

	claim.Redeemed = time.Now().Unix()
	return claims.Store(claim, false)
}

func CreateClaim(ctx context.Context, c *commands.CreateClaim, claims storage.ClaimStore) error {
	claim := &v1.Claim{
		Code:         c.Code,
		ResourceType: c.ResourceType,
		Created:      time.Now().Unix(),
		Redeemed:     0,
	}

	return claims.Store(claim, true)
}

func DeleteUser(ctx context.Context, c *commands.DeleteUser, users storage.UserStore) error {
	return users.Delete(c.Name)
}

func StoreUser(ctx context.Context, c *commands.StoreUser, users storage.UserStore) error {
	u := c.User
	if u.Name == "" {
		u.Name = slugify.Marshal(u.FullName)
	}

	return users.Store(u, c.New)
}

func FindUser(ctx context.Context, q *queries.FindUser, users storage.UserStore) (*v1.User, error) {
	return users.Find(q.Name)
}

func CountUsers(ctx context.Context, q *queries.CountUsers, users storage.UserStore) (int, error) {
	return users.Count(&storage.UserFilters{})
}

func SearchUsers(ctx context.Context, q *queries.SearchUsers, users storage.UserStore) ([]*v1.User, error) {
	return users.FindMany(&storage.UserFilters{})
}

func StorePost(ctx context.Context, c *commands.StorePost, posts storage.PostStore) error {
	p := c.Post
	if c.New {
		p.Created = time.Now().Unix()
	}

	if p.Name == "" {
		p.Name = slugify.Marshal(p.Title)
	}

	return posts.Store(p, c.New)
}

func DeletePost(ctx context.Context, c *commands.DeletePost, posts storage.PostStore) error {
	return posts.Delete(c.Name)
}

func SearchPosts(ctx context.Context, q *queries.SearchPosts, posts storage.PostStore) ([]*v1.Post, error) {
	return posts.FindMany(&storage.PostFilters{})
}

func FindPost(ctx context.Context, q *queries.FindPost, posts storage.PostStore) (*v1.Post, error) {
	return posts.Find(q.Name)
}
