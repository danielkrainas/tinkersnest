package version

import (
	"fmt"

	"github.com/danielkrainas/tinkersnest/cmd"
	"github.com/danielkrainas/tinkersnest/context"
)

func init() {
	cmd.Register("version", Info)
}

func run(ctx context.Context, args []string) error {
	fmt.Println("Tinkerctl v" + context.GetVersion(ctx))
	return nil
}

var (
	Info = &cmd.Info{
		Use:   "version",
		Short: "`version`",
		Long:  "`version`",
		Run:   cmd.ExecutorFunc(run),
	}
)
