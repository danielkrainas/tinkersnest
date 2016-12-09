package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/danielkrainas/tinkersnest/api/v1"
)

type UserAPI interface {
	SearchUsers() ([]*v1.User, error)
	CreateUser(user *v1.User) (*v1.User, error)
	UpdateUser(user *v1.User) (*v1.User, error)
	CreateUserWithClaim(user *v1.User, claim string) (*v1.User, error)
	GetUser(name string) (*v1.User, error)
	DeleteUser(name string) error
}

type usersAPI struct {
	*Client
}

func (c *Client) Users() UserAPI {
	return &usersAPI{c}
}

func (api *usersAPI) GetUser(name string) (*v1.User, error) {
	url, err := api.urls().BuildUserByName(name)
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

	p := &v1.User{}
	if err = json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	return p, nil
}

func (api *usersAPI) DeleteUser(name string) error {
	url, err := api.urls().BuildUserByName(name)
	if err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := api.do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (api *usersAPI) SearchUsers() ([]*v1.User, error) {
	url, err := api.urls().BuildUserRegistry()
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

	p := make([]*v1.User, 0)
	if err = json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	return p, nil
}

func (api *usersAPI) UpdateUser(user *v1.User) (*v1.User, error) {
	body, err := json.Marshal(&user)
	if err != nil {
		return nil, err
	}

	url, err := api.urls().BuildUserByName(user.Name)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
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

	u := &v1.User{}
	if err = json.Unmarshal(body, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (api *usersAPI) CreateUser(user *v1.User) (*v1.User, error) {
	body, err := json.Marshal(&user)
	if err != nil {
		return nil, err
	}

	url, err := api.urls().BuildUserRegistry()
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

	p := &v1.User{}
	if err = json.Unmarshal(body, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (api *usersAPI) CreateUserWithClaim(user *v1.User, claim string) (*v1.User, error) {
	body, err := json.Marshal(&user)
	if err != nil {
		return nil, err
	}

	url, err := api.urls().BuildUserRegistry()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	r.Header.Add("TINKERSNEST-CLAIM", claim)
	resp, err := api.do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	p := &v1.User{}
	if err = json.Unmarshal(body, p); err != nil {
		return nil, err
	}

	return p, nil
}
