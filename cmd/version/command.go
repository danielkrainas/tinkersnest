package version

import (
	"context"
	"fmt"

	"github.com/danielkrainas/tinkersnest/cmd"
	"github.com/danielkrainas/tinkersnest/context"
)

func init() {
	cmd.Register("version", Info)
}

func run(ctx context.Context, args []string) error {
	fmt.Println("TinkersNest v" + acontext.GetVersion(ctx))
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "version",
		Short: "show version information",
		Long:  "show version information",
		Run:   cmd.ExecutorFunc(run),
	}
)
