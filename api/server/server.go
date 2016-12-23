package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	cfg "github.com/danielkrainas/gobag/configuration"
	"github.com/danielkrainas/gobag/context"
	"github.com/danielkrainas/gobag/decouple/cqrs"
	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"github.com/danielkrainas/tinkersnest/api/server/handlers"
	"github.com/danielkrainas/tinkersnest/configuration"
	"github.com/danielkrainas/tinkersnest/setup"
	"github.com/danielkrainas/tinkersnest/storage"
)

type Server struct {
	context.Context
	config  *configuration.Config
	app     *handlers.App
	server  *http.Server
	query   *cqrs.QueryDispatcher
	command *cqrs.CommandDispatcher
	setup   *setup.SetupManager
}

func New(ctx context.Context, config *configuration.Config) (*Server, error) {
	ctx, err := configureLogging(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error configuring logging: %v", err)
	}

	log := acontext.GetLogger(ctx)
	log.Info("initializing server")

	setupManager := &setup.SetupManager{}

	storageDriver, err := storage.FromConfig(config)
	if err != nil {
		return nil, err
	}

	if err := storageDriver.Init(); err != nil {
		return nil, err
	}

	query := &cqrs.QueryDispatcher{
		Executors: []cqrs.QueryExecutor{
			setupManager,
			storageDriver.Query(),
		},
	}

	command := &cqrs.CommandDispatcher{
		Handlers: []cqrs.CommandHandler{
			setupManager,
			storageDriver.Command(),
		},
	}

	ctx = cqrs.WithCommandDispatch(ctx, command)
	ctx = cqrs.WithQueryDispatch(ctx, query)

	app, err := handlers.NewApp(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating server app: %v", err)
	}

	handler := alive("/", app)
	handler = panicHandler(handler)
	handler = contextHandler(app, handler)
	handler = loggingHandler(app, handler)

	n := negroni.New()

	n.Use(cors.New(cors.Options{
		AllowedOrigins:   config.HTTP.CORS.Origins,
		AllowedMethods:   config.HTTP.CORS.Methods,
		AllowCredentials: true,
		AllowedHeaders:   config.HTTP.CORS.Headers,
		Debug:            config.HTTP.CORS.Debug,
	}))

	n.UseHandler(handler)

	s := &Server{
		Context: ctx,
		app:     app,
		config:  config,
		query:   query,
		command: command,
		setup:   setupManager,
		server: &http.Server{
			Addr:    config.HTTP.Addr,
			Handler: n,
		},
	}

	log.Infof("using %q logging formatter", config.Log.Formatter)
	storage.LogSummary(ctx, config)

	if err := setupManager.Bootstrap(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

func (server *Server) ListenAndServe() error {
	config := server.config
	ln, err := net.Listen("tcp", config.HTTP.Addr)
	if err != nil {
		return err
	}

	acontext.GetLogger(server.app).Infof("listening on %v", ln.Addr())
	return server.server.Serve(ln)
}

func panicHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Panicf("%v", err)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

func alive(path string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == path {
			w.Header().Set("Cache-Control", "no-cache")
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func contextHandler(parent context.Context, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := acontext.DefaultContextManager.Context(parent, w, r)
		defer acontext.DefaultContextManager.Release(ctx)

		handler.ServeHTTP(w, r)
	})
}

func loggingHandler(parent context.Context, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := acontext.DefaultContextManager.Context(parent, w, r)
		acontext.GetRequestLogger(ctx).Info("request started")
		handler.ServeHTTP(w, r)
	})
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

		ctx = acontext.WithValues(ctx, config.Log.Fields)
		ctx = acontext.WithLogger(ctx, acontext.GetLogger(ctx, fields...))
	}

	ctx = acontext.WithLogger(ctx, acontext.GetLogger(ctx))
	return ctx, nil
}

func logLevel(level cfg.LogLevel) log.Level {
	l, err := log.ParseLevel(string(level))
	if err != nil {
		l = log.InfoLevel
		log.Warnf("error parsing level %q: %v, using %q", level, err, l)
	}

	return l
}
