package local

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/danielkrainas/tinkersnest/api/client"
)

const (
	TINKERCTL_HOME   = ".tinkerctl"
	AUTH_CONFIG_FILE = "auth.json"
)

type HostConfig struct {
	Host     string           `json:"-"`
	Username string           `json:"name"`
	Token    client.AuthToken `json:"token"`
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
		return "", err
	}

	return path.Join(u.HomeDir, TINKERCTL_HOME, AUTH_CONFIG_FILE), nil
}

func EnsureHomeExists() error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	homePath := path.Join(u.HomeDir, TINKERCTL_HOME)
	if _, err := os.Stat(homePath); err != nil {
		return os.MkdirAll(homePath, 0775)
	}

	return nil
}

func SaveAuthConfig(config AuthConfig) error {
	authConfigPath, err := getAuthConfigPath()
	if err != nil {
		return err
	}

	buf, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(authConfigPath, buf, 0644); err != nil {
		return err
	}

	return nil
}

func LoadAuthConfig() (AuthConfig, error) {
	authConfigPath, err := getAuthConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(authConfigPath); err != nil {
		return AuthConfig{}, nil
	}

	buf, err := ioutil.ReadFile(authConfigPath)
	if err != nil {
		return nil, err
	}

	config := AuthConfig{}
	if err := json.Unmarshal(buf, &config); err != nil {
		return nil, err
	}

	return config, nil
}
