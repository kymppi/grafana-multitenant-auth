package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
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

func NewRouter(lc fx.Lifecycle, logger *zap.Logger, config *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger(logger))

	// JWT verification key
	key, err := jwk.FromRaw([]byte(config.JWT_SECRET_KEY))

	if err != nil {
		logger.Fatal("Failed to create JWK key", zap.Error(err))
	}

	router.Get("/auth/", func(w http.ResponseWriter, r *http.Request) {
		// format Bearer <token>
		token := r.Header.Get("Authorization")[7:]

		parsed, err := jwt.ParseString(token, jwt.WithKey(jwa.HS256, key), jwt.WithIssuer(config.JWT_ALLOWED_ISSUER))

		if err != nil {
			logger.Error("Failed to parse token", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tenant, ok := parsed.Get("tenant")

		if !ok {
			logger.Error("Failed to get tenant from token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		logger.Info("Verification successful for a key with", zap.String("tenant", tenant.(string)), zap.String("issuer", parsed.Issuer()), zap.String("subject", parsed.Subject()))

		w.Header().Set("X-Scope-OrgID", tenant.(string))
		w.WriteHeader(http.StatusOK)
	})

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
				logger.Error("Failed to start server", zap.Error(err))
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
