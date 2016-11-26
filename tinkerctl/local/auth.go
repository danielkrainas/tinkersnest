package local

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
)

const (
	TINKERCTL_HOME   = ".tinkerctl"
	AUTH_CONFIG_FILE = "auth.json"
)

type HostConfig struct {
	Host     string `json:"-"`
	Username string `json:"name"`
	Token    string `json:"token"`
}

type AuthConfig map[string]*HostConfig

func (a AuthConfig) Exists(host string) bool {
	_, ok := a[host]
	return ok
}

func (a AuthConfig) Get(host string) *HostConfig {
	h, ok := a[host]
	if !ok {
		return nil
	}

	return h
}

func (a AuthConfig) Set(config *HostConfig) {
	a[config.Host] = config
}

func getAuthConfigPath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	return path.Join(u.HomeDir, TINKERCTL_HOME, AUTH_CONFIG_FILE)
}

func SaveAuthConfig(config AuthConfig) error {
	authConfigPath, err := getAuthConfigPath()
	if err != nil {
		return err
	}

	fr, err := os.Open(authConfigPath)

	return nil
}

func LoadAuthConfig() (AuthConfig, error) {
	authConfigPath, err := getAuthConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(authConfigPath); err != nil {
		return new(AuthConfig), nil
	}

	return new(AuthConfig), err
}
