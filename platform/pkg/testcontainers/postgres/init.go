package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startPostgresContainer(ctx context.Context, cfg *Config) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Name:     cfg.ContainerName,
		Image:    cfg.ImageName,
		Networks: []string{cfg.NetworkName},
		Env: map[string]string{
			postgresEnvUsernameKey: cfg.Username,
			postgresEnvPasswordKey: cfg.Password,
			postgresEnvDatabaseKey: cfg.Database,
		},
		WaitingFor:         wait.ForListeningPort(postgresPort + "/tcp").WithStartupTimeout(postgresStartupTimeout),
		HostConfigModifier: defaultHostConfig(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, errors.Errorf("failed to start postgres container: %v", err)
	}

	return container, nil
}

func getContainerHostPort(ctx context.Context, container testcontainers.Container) (string, string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", "", errors.Errorf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, postgresPort+"/tcp")
	if err != nil {
		return "", "", errors.Errorf("failed to get mapped port: %v", err)
	}

	return host, port.Port(), nil
}

func buildPostgresDSN(cfg *Config) string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Username,
		cfg.Password,
		cfg.SSLMode,
	)
}
