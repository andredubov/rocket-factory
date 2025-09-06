package iam

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type grpcClientEnvConfig struct {
	Host string `env:"IAM_GRPC_HOST,required"`
	Port string `env:"IAM_GRPC_PORT,required"`
}

type grpcClientConfig struct {
	raw grpcClientEnvConfig
}

func NewGRPCConfig() (*grpcClientConfig, error) {
	var raw grpcClientEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &grpcClientConfig{raw: raw}, nil
}

func (cfg *grpcClientConfig) Address() string {
	return net.JoinHostPort(
		cfg.raw.Host,
		cfg.raw.Port,
	)
}
