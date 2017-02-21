package create

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/tinkersnest/api/client"
	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/tinkerctl/resource"
)

func init() {
	cmd.Register("create", Info)
}

func run(ctx context.Context, args []string) error {
	claimCode, _ := ctx.Value("flags.claim").(string)
	resourcePath, ok := ctx.Value("flags.file").(string)
	if !ok {
		return errors.New("a resource file path is required")
	}

	res, err := resource.Load(resourcePath)
	if err != nil {
		return err
	}

	const ENDPOINT = "http://localhost:9240"

	c := client.New(ENDPOINT, http.DefaultClient)

	switch res.Type {
	case resource.Post:
		post, err := postFromSpec(res.Name, res.Spec)
		if err != nil {
			return err
		}

		if post, err = c.Blog().CreatePost(post); err != nil {
			return err
		}

		fmt.Printf("post %q was created!\n", res.Name)

	case resource.User:
		user, err := userFromSpec(res.Name, res.Spec)
		if err != nil {
			return err
		}

		if claimCode == "" {
			user, err = c.Users().CreateUser(user)
		} else {
			user, err = c.Users().CreateUserWithClaim(user, claimCode)
		}

		if err != nil {
			return err
		}

		fmt.Printf("user %q was created!\n", res.Name)

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
				Description: "resource spec file",
				Type:        cmd.FlagString,
			},
			{
				Short:       "c",
				Long:        "claim",
				Description: "claim code to redeem when creating a resource",
				Type:        cmd.FlagString,
			},
		},
	}
)

func userFromSpec(name string, spec map[string]interface{}) (*v1.User, error) {
	m, ok := spec["user"].(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("missing 'user' data in spec")
	}

	u := &v1.User{
		Name:     name,
		Email:    m["email"].(string),
		FullName: m["full_name"].(string),
		Password: m["password"].(string),
	}

	return u, nil
}

func postFromSpec(name string, spec map[string]interface{}) (*v1.Post, error) {
	m, ok := spec["post"].(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("missing 'post' data in spec")
	}

	p := &v1.Post{
		Name:    name,
		Title:   m["title"].(string),
		Publish: false,
		Content: make([]*v1.Content, 0),
		Tags:    make([]string, 0),
	}

	if author, ok := m["author"].(*v1.Author); ok {
		p.Author = author
	}

	if created, ok := m["created"].(int64); ok {
		p.Created = created
	}

	if publish, ok := m["publish"].(bool); ok {
		p.Publish = publish
	}

	if tags, ok := m["tags"].([]string); ok {
		p.Tags = tags
	}

	contents, ok := m["content"].([]interface{})
	if !ok {
		return nil, errors.New("missing 'content' in post spec")
	}

	for _, c := range contents {
		if cm, ok := c.(map[interface{}]interface{}); ok {
			c, err := getContent(cm)
			if err != nil {
				return nil, err
			}

			p.Content = append(p.Content, c)
		}
	}

	return p, nil
}

func getContent(spec map[interface{}]interface{}) (*v1.Content, error) {
	c := &v1.Content{}
	if t, ok := spec["type"].(string); !ok {
		return nil, errors.New("invalid or missing content 'type' in spec")
	} else {
		c.Type = t
	}

	if sdata, ok := spec["data"].(string); ok {
		c.Data = []byte(sdata)
	} else if src, ok := spec["src"].(string); ok {
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
