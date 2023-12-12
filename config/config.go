package config

import (
	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
)

type Config struct {
	Host   string `env:"HOST" envDefault:"0.0.0.0"`
	Port   uint16 `env:"PORT" envDefault:"3000"`
	GO_ENV string `env:"GO_ENV" envDefault:"production"`
}

func Parse(logger *zap.Logger) (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		logger.Error("Failed to parse config", zap.Error(err))
		return nil, err
	}

	logger.Info("Config loaded successfully")

	return &cfg, nil
}
