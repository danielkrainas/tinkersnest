package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Credentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AuthToken string

var (
	InvalidToken AuthToken
)

type AuthAPI interface {
	Login(c *Credentials) (AuthToken, error)
}

type authAPI struct {
	*Client
}

func (c *Client) Auth() AuthAPI {
	return &authAPI{c}
}

func (api *authAPI) Login(c *Credentials) (AuthToken, error) {
	body, err := json.Marshal(&c)
	if err != nil {
		return InvalidToken, err
	}

	url, err := api.urls().BuildAuth()
	if err != nil {
		return InvalidToken, err
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return InvalidToken, err
	}

	resp, err := api.do(r)
	if err != nil {
		return InvalidToken, err
	}

	defer resp.Body.Close()
	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return InvalidToken, err
	}

	return AuthToken(token), nil
}
