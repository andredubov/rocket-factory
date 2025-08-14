package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type paymentClientEnvConfig struct {
	Host string `env:"PAYMENT_GRPC_HOST,required"`
	Port string `env:"PAYMENT_GRPC_PORT,required"`
}

type paymentGrpcConfig struct {
	raw paymentClientEnvConfig
}

func NewGRPCConfig() (*paymentGrpcConfig, error) {
	var raw paymentClientEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &paymentGrpcConfig{raw: raw}, nil
}

func (cfg *paymentGrpcConfig) Address() string {
	return net.JoinHostPort(
		cfg.raw.Host,
		cfg.raw.Port,
	)
}
