package resource

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/danielkrainas/tinkersnest/configuration"
)

type v1_0Resource Resource

func Parse(rd io.Reader) (*Resource, error) {
	in, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	p := configuration.NewParser("tinkerctl", []configuration.VersionedParseInfo{
		{
			Version: configuration.MajorMinorVersion(1, 0),
			ParseAs: reflect.TypeOf(v1_0Resource{}),
			ConversionFunc: func(c interface{}) (interface{}, error) {
				if v1_0, ok := c.(*v1_0Resource); ok {
					return (*Resource)(v1_0), nil
				}

				return nil, fmt.Errorf("Expected *v1_0Resource, received %#v", c)
			},
		},
	})

	res := new(Resource)
	err = p.Parse(in, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
