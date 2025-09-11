package env

import (
	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level              string `env:"LOGGER_LEVEL,required"`
	AsJson             bool   `env:"LOGGER_AS_JSON,required"`
	EnableOTLP         bool   `env:"LOGGER_ENABLE_OTLP,required"`
	OTLPEndpoint       string `env:"LOGGER_OTLP_ENDPOINT,required"`
	ServiceName        string `env:"LOGGER_SERVICE_NAME,required"`
	ServiceEnvironment string `env:"LOGGER_SERVICE_ENVIRONMENT,required"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &loggerConfig{raw: raw}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *loggerConfig) AsJson() bool {
	return cfg.raw.AsJson
}

func (l *loggerConfig) EnableOTLP() bool {
	return l.raw.EnableOTLP
}

func (l *loggerConfig) OTLPEndpoint() string {
	return l.raw.OTLPEndpoint
}

func (l *loggerConfig) ServiceName() string {
	return l.raw.ServiceName
}

func (l *loggerConfig) ServiceEnvironment() string {
	return l.raw.ServiceEnvironment
}
