package get

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/danielkrainas/tinkersnest/api/client"
	"github.com/danielkrainas/tinkersnest/cmd"
)

func init() {
	cmd.Register("get", Info)
}

func run(ctx context.Context, args []string) error {
	if len(args) < 1 || args[0] == "" {
		return errors.New("you must specify a resource type")
	}

	const ENDPOINT = "http://localhost:9240"

	c := client.New(ENDPOINT, http.DefaultClient)

	switch args[0] {
	case "posts":
		posts, err := c.Blog().SearchPosts()
		if err != nil {
			return err
		}

		fmt.Printf("%10s | %-20s \n", "PUBLISHED", "NAME")
		for _, post := range posts {
			fmt.Printf("%10s | %-20s \n", yesNoBool(post.Publish), post.Name)
		}

		fmt.Println("")

	default:
		return fmt.Errorf("resource type %q unsupported", args[0])
	}

	return nil
}

var (
	Info = &cmd.Info{
		Use:   "get <resource_type>",
		Short: "list a type of resources on the server",
		Long:  "list a type of resources on the server",
		Run:   cmd.ExecutorFunc(run),
	}
)

func yesNoBool(b bool) string {
	if b {
		return "yes"
	}

	return "no"
}
