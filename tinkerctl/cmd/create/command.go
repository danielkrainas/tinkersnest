package create

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/danielkrainas/tinkersnest/api/client"
	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/cmd"
	"github.com/danielkrainas/tinkersnest/tinkerctl/spec"
)

func init() {
	cmd.Register("create", Info)
}

func run(ctx context.Context, args []string) error {
	specPath, ok := ctx.Value("flags.file").(string)
	if !ok {
		return errors.New("an object spec file path is required")
	}

	res, err := spec.Load(specPath)
	if err != nil {
		return err
	}

	const ENDPOINT = "http://localhost:9240"

	c := client.New(ENDPOINT, http.DefaultClient)

	switch res.Type {
	case spec.Post:
		post, err := postFromSpec(res.Name, res.Spec)
		if err != nil {
			return err
		}

		if post, err = c.Blog().CreatePost(post); err != nil {
			return err
		}

		fmt.Printf("post %q was created!\n", res.Name)

	default:
		return fmt.Errorf("resource type %q unsupported", res.Type)
	}

	return nil
}

var (
	Info = &cmd.Info{
		Use:   "create",
		Short: "create a resource on the server",
		Long:  "create a resource on the server",
		Run:   cmd.ExecutorFunc(run),
		Flags: []*cmd.Flag{
			{
				Short:       "f",
				Long:        "file",
				Description: "resource meta file",
				Type:        cmd.FlagString,
			},
		},
	}
)

func postFromSpec(name string, meta map[string]interface{}) (*v1.Post, error) {
	m, ok := meta["post"].(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("missing 'post' data in metadata")
	}

	p := &v1.Post{
		Name:    name,
		Title:   m["title"].(string),
		Publish: false,
		Content: make([]*v1.Content, 0),
	}

	if created, ok := m["created"].(int64); ok {
		p.Created = created
	}

	if publish, ok := m["publish"].(bool); ok {
		p.Publish = publish
	}

	contents, ok := m["content"].([]interface{})
	if !ok {
		return nil, errors.New("missing 'content' in post metadata")
	}

	for _, c := range contents {
		if cm, ok := c.(map[string]interface{}); ok {
			c, err := getContent(cm)
			if err != nil {
				return nil, err
			}

			p.Content = append(p.Content, c)
		}
	}

	return p, nil
}

func getContent(meta map[string]interface{}) (*v1.Content, error) {
	c := &v1.Content{}
	if t, ok := meta["type"].(string); !ok {
		return nil, errors.New("invalid or missing content 'type' in metadata")
	} else {
		c.Type = t
	}

	if sdata, ok := meta["data"].(string); ok {
		c.Data = []byte(sdata)
	} else if src, ok := meta["src"].(string); ok {
		data, err := ioutil.ReadFile(src)
		if err != nil {
			return nil, err
		}

		c.Data = data
	}

	if c.Data == nil {
		return nil, errors.New("content does not have any data associated")
	}

	return c, nil
}
