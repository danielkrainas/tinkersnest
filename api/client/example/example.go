package main

import (
	"fmt"
	"net/http"

	"github.com/danielkrainas/tinkersnest/api/client"
	"github.com/danielkrainas/tinkersnest/api/v1"
)

func main() {
	const ENDPOINT = "http://localhost:9240"

	// Create a new client
	c := client.New(ENDPOINT, http.DefaultClient)
	fmt.Printf("created new client to %q\n", ENDPOINT)

	// Check V1 endpoint is good and healthy
	//=====================================
	err := c.Ping()
	if err != nil {
		panic("error sending ping")
	}

	fmt.Println("sent ping")

	// Create a new blog post
	//=====================================
	post, err := c.Blog().CreatePost(&v1.CreatePostRequest{
		Title: "Example Post",
		Body:  "This was an example post using the tinkersnest/api/client package.",
	})

	if err != nil {
		panic("error creating post")
	}

	fmt.Printf("created %q post.\n", post.Title)

	// Search the blog posts
	//=====================================
	posts, err := c.Blog().SearchPosts()
	if err != nil {
		panic("error searching posts")
	}

	fmt.Printf("found %d posts\n", len(posts))
}
