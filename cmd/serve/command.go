package serve

import (
	"context"

	"github.com/danielkrainas/gobag/cmd"

	"github.com/danielkrainas/tinkersnest/api/server"
	"github.com/danielkrainas/tinkersnest/configuration"
)

func init() {
	cmd.Register("serve", Info)
}

func run(ctx context.Context, args []string) error {
	config, err := configuration.Resolve(args)
	if err != nil {
		return err
	}

	s, err := server.New(ctx, config)
	if err != nil {
		return err
	}

	return s.ListenAndServe()
}

var (
	Info = &cmd.Info{
		Use:   "serve",
		Short: "run the api server",
		Long:  "run the api server",
		Run:   cmd.ExecutorFunc(run),
	}
)
