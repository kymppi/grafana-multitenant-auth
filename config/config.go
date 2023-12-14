package config

import (
	"errors"

	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
)

type Config struct {
	Host               string `env:"HOST" envDefault:"0.0.0.0"`
	Port               uint16 `env:"PORT" envDefault:"3000"`
	GO_ENV             string `env:"GO_ENV" envDefault:"production"`
	JWT_ALLOWED_ISSUER string `env:"JWT_ALLOWED_ISSUER" envDefault:""`
	JWT_SECRET_KEY     string `env:"JWT_SECRET_KEY" envDefault:""`
}

func Parse(logger *zap.Logger) (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		logger.Error("Failed to parse config", zap.Error(err))
		return nil, err
	}

	logger.Info("Config loaded successfully")

	// TODO: validate config
	if cfg.JWT_ALLOWED_ISSUER == "" {
		logger.Fatal("JWT_ALLOWED_ISSUER is required")
		return nil, errors.New("JWT_ALLOWED_ISSUER is required")
	}

	if cfg.JWT_SECRET_KEY == "" {
		logger.Fatal("JWT_SECRET_KEY is required")
		return nil, errors.New("JWT_SECRET_KEY is required")
	}

	return &cfg, nil
}
