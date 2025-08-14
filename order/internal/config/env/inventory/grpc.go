package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryClientEnvConfig struct {
	Host string `env:"INVENTORY_GRPC_HOST,required"`
	Port string `env:"INVENTORY_GRPC_PORT,required"`
}

type inventoryGrpcConfig struct {
	raw inventoryClientEnvConfig
}

func NewGRPCConfig() (*inventoryGrpcConfig, error) {
	var raw inventoryClientEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &inventoryGrpcConfig{raw: raw}, nil
}

func (cfg *inventoryGrpcConfig) Address() string {
	return net.JoinHostPort(
		cfg.raw.Host,
		cfg.raw.Port,
	)
}
