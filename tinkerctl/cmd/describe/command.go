package describe

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/danielkrainas/tinkersnest/api/client"
	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/cmd"
)

func init() {
	cmd.Register("describe", Info)
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
	case "post":
		post, err := c.Blog().GetPost(name)
		if err != nil {
			return err
		}

		describePost(post)

	default:
		return fmt.Errorf("resource type %q unsupported", args[0])
	}

	return nil
}

var (
	Info = &cmd.Info{
		Use:   "describe <resource_type> <name>",
		Short: "show details about a resource",
		Long:  "show details about a resource",
		Run:   cmd.ExecutorFunc(run),
	}
)

func describePost(p *v1.Post) {
	fmt.Println("[metadata]")
	fmt.Printf("name:  %s\n", p.Name)
	fmt.Printf("title:  %s\n", p.Title)
	fmt.Printf("created:  %d\n", p.Created)
	fmt.Printf("publish: %s\n", yesNoBool(p.Publish))
	for i, c := range p.Content {
		fmt.Printf("[content#%d %s]\n", i, c.Type)
		fmt.Println(string(c.Data))
	}

	fmt.Print("\n")
}

func yesNoBool(b bool) string {
	if b {
		return "yes"
	}

	return "no"
}
