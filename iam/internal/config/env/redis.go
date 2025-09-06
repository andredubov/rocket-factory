package env

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type redisEnv struct {
	Host              string        `env:"REDIS_HOST,required"`
	Port              int           `env:"REDIS_PORT,required"`
	ConnectionTimeout time.Duration `env:"REDIS_CONNECTION_TIMEOUT,required"`
	MaxIdle           int           `env:"REDIS_MAX_IDLE,required"`
	IdleTimeout       time.Duration `env:"REDIS_IDLE_TIMEOUT,required"`
	SessionTTL        time.Duration `env:"SESSION_TTL,required"`
}

type redisConfig struct {
	raw redisEnv
}

func (r *redisConfig) Address() string {
	return fmt.Sprintf("%s:%d", r.raw.Host, r.raw.Port)
}

func (r *redisConfig) ConnectionTimeout() time.Duration {
	return r.raw.ConnectionTimeout
}

func (r *redisConfig) MaxIdle() int {
	return r.raw.MaxIdle
}

func (r *redisConfig) IdleTimeout() time.Duration {
	return r.raw.IdleTimeout
}

func (r *redisConfig) CacheTTL() time.Duration {
	return r.raw.SessionTTL
}

func NewRedisConfig() (*redisConfig, error) {
	var raw redisEnv
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &redisConfig{raw: raw}, nil
}
