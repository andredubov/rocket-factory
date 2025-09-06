package redis

import (
	"context"
	"net"
	"time"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/andredubov/rocket-factory/platform/pkg/cache"
	rediscache "github.com/andredubov/rocket-factory/platform/pkg/cache/redis"
)

func createRedisPool(cfg *Config) *redigo.Pool {
	address := net.JoinHostPort(cfg.Host, cfg.Port)

	connPool := &redigo.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   20,
		IdleTimeout: cfg.IdleTimeout,
		Wait:        true,
		DialContext: func(ctx context.Context) (redigo.Conn, error) {
			conn, err := redigo.DialContext(ctx, "tcp", address,
				redigo.DialConnectTimeout(cfg.ConnectionTimeout),
				redigo.DialReadTimeout(cfg.ConnectionTimeout),
				redigo.DialWriteTimeout(cfg.ConnectionTimeout),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return connPool
}

func getRedisClient(cfg *Config) cache.RedisClient {
	pool := createRedisPool(cfg)

	return rediscache.NewClient(
		pool,
		cfg.Logger,
		cfg.ConnectionTimeout,
	)
}
