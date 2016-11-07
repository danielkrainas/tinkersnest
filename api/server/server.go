package server

import (
	"fmt"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/rs/cors"
	"github.com/urfave/negroni"

	"github.com/danielkrainas/tinkersnest/api/server/handlers"
	"github.com/danielkrainas/tinkersnest/configuration"
	"github.com/danielkrainas/tinkersnest/context"
)

type Server struct {
	config *configuration.Config
	app    *handlers.App
	server *http.Server
}

func New(ctx context.Context, config *configuration.Config) (*Server, error) {
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
		app:    app,
		config: config,
		server: &http.Server{
			Addr:    config.HTTP.Addr,
			Handler: n,
		},
	}

	log.Infof("using %q logging formatter", config.Log.Formatter)
	log.Infof("using %q containers driver", config.Containers.Type())
	log.Infof("using %q storage driver", config.Storage.Type())
	if !config.HTTP.Enabled {
		log.Info("http api disabled")
	}

	return s, nil
}

func (server *Server) ListenAndServe() error {
	config := server.config
	ln, err := net.Listen("tcp", config.HTTP.Addr)
	if err != nil {
		return err
	}

	context.GetLogger(server.app).Infof("listening on %v", ln.Addr())
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
		ctx := context.DefaultContextManager.Context(parent, w, r)
		defer context.DefaultContextManager.Release(ctx)

		handler.ServeHTTP(w, r)
	})
}

func loggingHandler(parent context.Context, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.DefaultContextManager.Context(parent, w, r)
		context.GetRequestLogger(ctx).Info("request started")
		handler.ServeHTTP(w, r)
	})
}