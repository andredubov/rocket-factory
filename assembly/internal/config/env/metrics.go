package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type metricsEnvConfig struct {
	CollectorEndpoint string        `env:"METRICS_OTEL_COLLECTOR_ENDPOINT,required"`
	CollectorInterval time.Duration `env:"METRICS_OTEL_COLLECTOR_INTERVAL,required"`
	CollectorTimeout  time.Duration `env:"METRICS_OTEL_COLLECTOR_TIMEOUT"`
}

type metricsConfig struct {
	raw metricsEnvConfig
}

func NewMetricsConfig() (*metricsConfig, error) {
	var raw metricsEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &metricsConfig{raw: raw}, nil
}

func (m *metricsConfig) CollectorEndpoint() string { return m.raw.CollectorEndpoint }

func (m *metricsConfig) CollectorInterval() time.Duration {
	return m.raw.CollectorInterval
}

func (m *metricsConfig) CollectorTimeout() time.Duration {
	return m.raw.CollectorTimeout
}
