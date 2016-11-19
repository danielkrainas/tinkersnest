package resource

import (
	"fmt"
	"os"
)

type Resource struct {
	Name string                 `yaml:"name"`
	Type ResourceType           `yaml:"type"`
	Spec map[string]interface{} `yaml:"spec"`
}

type ResourceType string

var (
	Post ResourceType = "Post"
)

func Load(resourcePath string) (*Resource, error) {
	if resourcePath == "" {
		return nil, fmt.Errorf("Resource path not specified")
	}

	fp, err := os.Open(resourcePath)
	if err != nil {
		return nil, err
	}

	defer fp.Close()
	res, err := Parse(fp)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %v", resourcePath, err)
	}

	return res, nil
}
