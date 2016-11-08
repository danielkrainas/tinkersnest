package post

import (
	"github.com/danielkrainas/tinkersnest/cmd"
	"github.com/danielkrainas/tinkersnest/context"
)

func init() {
	cmd.Register("post", Info)
}

func run(ctx context.Context, args []string) error {
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "post",
		Short: "`post`",
		Long:  "`post`",
		Run:   cmd.ExecutorFunc(run),
	}
)
