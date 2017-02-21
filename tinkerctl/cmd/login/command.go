package describe

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/danielkrainas/gobag/cmd"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/danielkrainas/tinkersnest/api/client"
	"github.com/danielkrainas/tinkersnest/tinkerctl/local"
)

func init() {
	cmd.Register("login", Info)
}

func run(ctx context.Context, args []string) error {
	if len(args) < 1 || args[0] == "" {
		return errors.New("you must specify an endpoint url")
	}

	endpoint := args[0]

	if err := local.EnsureHomeExists(); err != nil {
		return err
	}

	config, err := local.LoadAuthConfig()
	if err != nil {
		return err
	}

	c := client.New(endpoint, http.DefaultClient)

	r := bufio.NewReader(os.Stdin)
	fmt.Print("username: ")
	username, err := r.ReadString('\n')
	if err != nil {
		return err
	}

	fmt.Printf("password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}

	username = strings.TrimSpace(username)
	passwordStr := strings.TrimSpace(string(password))

	token, err := c.Auth().Login(username, passwordStr)
	if err != nil {
		return err
	}

	config.Set(&local.HostConfig{
		Host:     endpoint,
		Username: username,
		Token:    token,
	})

	if err = local.SaveAuthConfig(config); err != nil {
		return err
	}

	return nil
}

var (
	Info = &cmd.Info{
		Use:   "login <url>",
		Short: "authenticate with a tinkersnest host",
		Long:  "authenticate with a tinkersnest host",
		Run:   cmd.ExecutorFunc(run),
	}
)
