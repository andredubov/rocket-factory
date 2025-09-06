package redis

import (
	"context"
	"errors"
	"time"

	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Config struct {
	NetworkName       string
	ContainerName     string
	ImageName         string
	ConnectionTimeout time.Duration
	IdleTimeout       time.Duration
	MaxIdle           int
	Logger            Logger
	Host              string
	Port              string
}

func validateConfig(cfg *Config) error {
	if cfg.MaxIdle <= 0 {
		return errors.New("MaxIdle must be positive")
	}
	if cfg.ConnectionTimeout <= 0 {
		return errors.New("ConnectionTimeout must be positive")
	}
	if cfg.IdleTimeout <= 0 {
		return errors.New("IdleTimeout must be positive")
	}
	return nil
}

func buildConfig(opts ...Option) (*Config, error) {
	cfg := &Config{
		NetworkName:       "test-network",
		ContainerName:     "test-redis-container",
		ImageName:         "redis:7.2.5-alpine3.20",
		ConnectionTimeout: 10 * time.Second,
		IdleTimeout:       10 * time.Second,
		MaxIdle:           10,
		Logger:            &logger.NoopLogger{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func defaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = true
	}
}
