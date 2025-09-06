package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamClientEnvConfig struct {
	Host string `env:"IAM_GRPC_HOST,required"`
	Port string `env:"IAM_GRPC_PORT,required"`
}

type iamGrpcConfig struct {
	raw iamClientEnvConfig
}

func NewGRPCConfig() (*iamGrpcConfig, error) {
	var raw iamClientEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &iamGrpcConfig{raw: raw}, nil
}

func (cfg *iamGrpcConfig) Address() string {
	return net.JoinHostPort(
		cfg.raw.Host,
		cfg.raw.Port,
	)
}
