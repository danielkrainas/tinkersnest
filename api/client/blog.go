package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/danielkrainas/tinkersnest/api/v1"
)

type BlogAPI interface {
	SearchPosts() ([]*v1.Post, error)
	CreatePost(create *v1.CreatePostRequest) (*v1.Post, error)
}

type blogAPI struct {
	*Client
}

func (c *Client) Blog() BlogAPI {
	return &blogAPI{c}
}

func (api *blogAPI) SearchPosts() ([]*v1.Post, error) {
	url, err := api.urls().BuildBlog()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	p := make([]*v1.Post, 0)
	if err = json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	return p, nil
}

func (api *blogAPI) CreatePost(create *v1.CreatePostRequest) (*v1.Post, error) {
	body, err := json.Marshal(&create)
	if err != nil {
		return nil, err
	}

	url, err := api.urls().BuildBlog()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := api.do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	p := &v1.Post{}
	if err = json.Unmarshal(body, p); err != nil {
		return nil, err
	}

	return p, nil
}
