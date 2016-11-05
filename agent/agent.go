package agent

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/danielkrainas/tinkersnest/api/server"
	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/configuration"
	"github.com/danielkrainas/tinkersnest/containers"
	containersFactory "github.com/danielkrainas/tinkersnest/containers/factory"
	"github.com/danielkrainas/tinkersnest/context"
	"github.com/danielkrainas/tinkersnest/hooks"
	"github.com/danielkrainas/tinkersnest/storage"
	storageDriverFactory "github.com/danielkrainas/tinkersnest/storage/factory"
)

type Agent struct {
	context.Context

	config *configuration.Config

	storage storage.Driver

	server *server.Server
}

func (agent *Agent) Run() error {
	context.GetLogger(agent).Info("starting agent")
	defer context.GetLogger(agent).Info("shutting down agent")

	go agent.server.ListenAndServe()
	agent.ProcessEvents()
	return nil
}

func (agent *Agent) getHostInfo() *v1.HostInfo {
	hostname, _ := os.Hostname()
	return &v1.HostInfo{
		Hostname: hostname,
	}
}

func New(ctx context.Context, config *configuration.Config) (*Agent, error) {
	ctx, err := configureLogging(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error configuring logging: %v", err)
	}

	log := context.GetLogger(ctx)
	log.Info("initializing agent")

	ctx, storageDriver, err := configureStorage(ctx, config)
	if err != nil {
		return nil, err
	}

	server, err := server.New(ctx, config)
	if err != nil {
		return nil, err
	}

	log.Infof("using %q logging formatter", config.Log.Formatter)
	log.Infof("using %q storage driver", config.Storage.Type())
	if !config.HTTP.Enabled {
		log.Info("http api disabled")
	}

	return &Agent{
		Context:    ctx,
		config:     config,
		containers: containersDriver,
		storage:    storageDriver,
		server:     server,
		hookFilter: &hooks.CriteriaFilter{},
		shooter:    &hooks.LiveShooter{http.DefaultClient},
	}, nil
}

func configureStorage(ctx context.Context, config *configuration.Config) (context.Context, storage.Driver, error) {
	storageParams := config.Storage.Parameters()
	if storageParams == nil {
		storageParams = make(configuration.Parameters)
	}

	storageDriver, err := storageDriverFactory.Create(config.Storage.Type(), storageParams)
	if err != nil {
		return ctx, nil, err
	}

	if err := storageDriver.Init(); err != nil {
		return ctx, nil, err
	}

	return storage.ForContext(ctx, storageDriver), storageDriver, nil
}

func configureLogging(ctx context.Context, config *configuration.Config) (context.Context, error) {
	log.SetLevel(logLevel(config.Log.Level))
	formatter := config.Log.Formatter
	if formatter == "" {
		formatter = "text"
	}

	switch formatter {
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})

	case "text":
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
		})

	default:
		if config.Log.Formatter != "" {
			return ctx, fmt.Errorf("unsupported log formatter: %q", config.Log.Formatter)
		}
	}

	if len(config.Log.Fields) > 0 {
		var fields []interface{}
		for k := range config.Log.Fields {
			fields = append(fields, k)
		}

		ctx = context.WithValues(ctx, config.Log.Fields)
		ctx = context.WithLogger(ctx, context.GetLogger(ctx, fields...))
	}

	ctx = context.WithLogger(ctx, context.GetLogger(ctx))
	return ctx, nil
}

func logLevel(level configuration.LogLevel) log.Level {
	l, err := log.ParseLevel(string(level))
	if err != nil {
		l = log.InfoLevel
		log.Warnf("error parsing level %q: %v, using %q", level, err, l)
	}

	return l
}
