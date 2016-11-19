package post

import (
	"context"

	"github.com/danielkrainas/tinkersnest/cmd"
)

func init() {
	cmd.Register("create", Info)
}

func run(ctx context.Context, args []string) error {
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "create",
		Short: "create an object on the server",
		Long:  "create an object on the server",
		Run:   cmd.ExecutorFunc(run),
		Flags: []*cmd.Flag{
			{
				Short:       "f",
				Long:        "file",
				Description: "object spec file",
				Type:        cmd.FlagString,
			},
		},
	}
)
