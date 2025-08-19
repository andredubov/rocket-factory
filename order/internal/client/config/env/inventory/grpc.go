package env

import (
	"fmt"
	"net"
	"os"

	"github.com/andredubov/rocket-factory/order/internal/config"
)

const (
	grpcHostEnvName = "INVENTORY_GRPC_HOST"
	grpcPortEnvName = "INVENTORY_GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

// NewGRPCConfig returns an instance of grpcConfig struct
func NewGRPCConfig() (config.GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("%s", "inventory grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("%s", "inventory grpc port not found")
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

// Address returns grpc server address
func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
