package delete

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/danielkrainas/tinkersnest/api/client"
	"github.com/danielkrainas/tinkersnest/cmd"
)

func init() {
	cmd.Register("delete", Info)
}

func run(ctx context.Context, args []string) error {
	if len(args) < 1 || args[0] == "" {
		return errors.New("you must specify a resource type")
	} else if len(args) < 2 || args[1] == "" {
		return errors.New("you must specify the name of a resource")
	}

	const ENDPOINT = "http://localhost:9240"

	c := client.New(ENDPOINT, http.DefaultClient)

	name := args[1]
	switch args[0] {
	case "user":
		user, err := c.Users().DeleteUser(name)
		if err != nil {
			return err
		}

	case "post":
		post, err := c.Blog().DeletePost(name)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("resource type %q unsupported", args[0])
	}

	return nil
}

var (
	Info = &cmd.Info{
		Use:   "delete <resource_type> <name>",
		Short: "delete a resource on the server",
		Long:  "delete a resource on the server",
		Run:   cmd.ExecutorFunc(run),
	}
)
