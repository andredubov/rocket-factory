package redis

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/platform/pkg/cache"
)

const (
	redisPort           = "6379"
	redisStartupTimeout = 1 * time.Minute

	redisEnvConnTimeoutKey = "REDIS_CONNECTION_TIMEOUT"
	redisEnvIdleTimeoutKey = "REDIS_IDLE_TIMEOUT"
	redisEnvMaxIdleKey     = "REDIS_MAX_IDLE"
)

type Container struct {
	container   testcontainers.Container
	redisClient cache.RedisClient
	cfg         *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg, err := buildConfig(opts...)
	if err != nil {
		return nil, err
	}

	container, err := startRedisContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	success := false
	defer func() {
		if !success {
			if err = container.Terminate(ctx); err != nil {
				cfg.Logger.Error(ctx, "failed to terminate postgres container", zap.Error(err))
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)
	if err != nil {
		return nil, err
	}

	redisClient := getRedisClient(cfg)

	cfg.Logger.Info(ctx, "Redis container started")
	success = true

	return &Container{
		container:   container,
		redisClient: redisClient,
		cfg:         cfg,
	}, nil
}

func (c *Container) RedisClient() cache.RedisClient {
	return c.redisClient
}

func (c *Container) Config() *Config {
	return c.cfg
}

func (c *Container) Terminate(ctx context.Context) error {
	if client, ok := c.redisClient.(interface{ Close() error }); ok {
		if err := client.Close(); err != nil {
			c.cfg.Logger.Error(ctx, "failed to close redis client", zap.Error(err))
		}
	}

	if err := c.container.Terminate(ctx); err != nil {
		c.cfg.Logger.Error(ctx, "failed to terminate postgres container", zap.Error(err))
	}

	c.cfg.Logger.Info(ctx, "Redis container terminated")

	return nil
}
