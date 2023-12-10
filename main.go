package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/kymppi/grafana-multitenant-auth-server/config"
	"github.com/kymppi/grafana-multitenant-auth-server/middleware"
)

func NewLogger() *zap.Logger {
	var logger *zap.Logger

	// use os.Getenv here since config has not been loaded yet
	// this is because config uses the logger
	if os.Getenv("GO_ENV") != "development" {
		logger = zap.Must(zap.NewProduction())
	} else {
		logger = zap.Must(zap.NewDevelopment())
	}

	defer logger.Sync()

	return logger
}

func NewRouter(lc fx.Lifecycle, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger(logger))

	return router
}

func NewHTTPServer(lc fx.Lifecycle, router *chi.Mux, logger *zap.Logger, config *config.Config) *http.Server {
	server := &http.Server{
		Addr:    config.Host + ":" + fmt.Sprint(config.Port),
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", server.Addr)

			if err != nil {
				return err
			}

			logger.Info("Starting HTTP Server at", zap.String("addr", server.Addr))

			go server.Serve(listener)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return server
}

func main() {
	fx.New(
		fx.NopLogger, // disable fx logging
		fx.Provide(
			NewLogger,
			config.Parse,
			NewRouter,
			NewHTTPServer,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
