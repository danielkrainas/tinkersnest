package serve

import (
	"github.com/danielkrainas/tinkersnest/agent"
	"github.com/danielkrainas/tinkersnest/cmd"
	"github.com/danielkrainas/tinkersnest/configuration"
	"github.com/danielkrainas/tinkersnest/context"
)

func init() {
	cmd.Register("serve", Info)
}

func run(ctx context.Context, args []string) error {
	config, err := configuration.Resolve(args)
	if err != nil {
		return err
	}

	agent, err := agent.New(ctx, config)
	if err != nil {
		return err
	}

	return agent.Run()
}

var (
	Info = &cmd.Info{
		Use:   "serve",
		Short: "`serve`",
		Long:  "`serve`",
		Run:   cmd.ExecutorFunc(run),
	}
)
