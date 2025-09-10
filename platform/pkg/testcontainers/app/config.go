package app

import (
	"errors"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

type Config struct {
	Name          string
	DockerfileDir string
	Dockerfile    string
	Port          string
	Env           map[string]string
	Networks      []string
	LogOutput     io.Writer
	StartupWait   wait.Strategy
	Logger        Logger
}

func validateConfig(cfg *Config) error {
	if cfg.Name == "" {
		return errors.New("name is required")
	}
	if cfg.DockerfileDir == "" {
		return errors.New("dockerfile directory is required")
	}
	if cfg.Dockerfile == "" {
		return errors.New("dockerfile is required")
	}
	return nil
}

func buildConfig(opts ...Option) (*Config, error) {
	cfg := &Config{
		Name:          defaultAppName,
		Port:          defaultAppPort,
		Dockerfile:    "Dockerfile",
		DockerfileDir: ".",
		LogOutput:     io.Discard,
		StartupWait:   wait.ForListeningPort(defaultAppPort + "/tcp").WithStartupTimeout(defaultStartupTimeout),
		Env:           make(map[string]string),
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

func DefaultHostConfig() func(hc *container.HostConfig) {
	return func(hc *container.HostConfig) {
		hc.AutoRemove = true
	}
}
