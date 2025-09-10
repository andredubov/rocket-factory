package postgres

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Config struct {
	NetworkName   string
	ContainerName string
	ImageName     string
	Database      string
	Username      string
	Password      string
	SSLMode       string
	MigrationDir  string
	Logger        Logger
	Host          string
	Port          string
}

func validateConfig(cfg *Config) error {
	if cfg.Username == "" {
		return errors.New("username is required")
	}
	if cfg.Password == "" {
		return errors.New("password is required")
	}
	if cfg.Database == "" {
		return errors.New("database is required")
	}
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	return nil
}

func buildConfig(opts ...Option) (*Config, error) {
	cfg := &Config{
		NetworkName:   "test-network",
		ContainerName: "test-postgres-container",
		ImageName:     "postgres:17.0-alpine3.20",
		Database:      "demo",
		Username:      "demo",
		Password:      "demo",
		SSLMode:       "disable",
		MigrationDir:  "migrations",
		Logger:        &logger.NoopLogger{},
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
