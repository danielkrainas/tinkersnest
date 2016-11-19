package spec

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/danielkrainas/tinkersnest/configuration"
)

type v1_0Spec Spec

func Parse(rd io.Reader) (*Spec, error) {
	in, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	p := configuration.NewParser("tinkerctl", []configuration.VersionedParseInfo{
		{
			Version: configuration.MajorMinorVersion(1, 0),
			ParseAs: reflect.TypeOf(v1_0Spec{}),
			ConversionFunc: func(c interface{}) (interface{}, error) {
				if v1_0, ok := c.(*v1_0Spec); ok {
					return (*Spec)(v1_0), nil
				}

				return nil, fmt.Errorf("Expected *v1_0Spec, received %#v", c)
			},
		},
	})

	spec := new(Spec)
	err = p.Parse(in, spec)
	if err != nil {
		return nil, err
	}

	return spec, nil
}
