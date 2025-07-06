package app

import (
	"context"
	"log"

	"github.com/andredubov/golibs/pkg/config"
	"github.com/andredubov/golibs/pkg/config/env"

	server "github.com/andredubov/rocket-factory/inventory/internal/api/v1/inventory"
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/memory"
)

// serviceProvider implements the dependency container pattern
// It provides lazy initialization of application components
type serviceProvider struct {
	inventoryRepository  repository.Inventory   // Inventory data access layer
	grpcConfig           config.GRPCConfig      // GRPC server configuration
	serverImplementation *server.Implementation // GRPC service implementation
}

// newServiceProvider creates a new service provider instance
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// GRPCConfig loads GRPC configuration from environment variables
// Implements singleton pattern - initializes config only once
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// InventoryRepository provides access to inventory data
// Uses in-memory implementation and singleton pattern
func (s *serviceProvider) InventoryRepository(ctx context.Context) repository.Inventory {
	if s.inventoryRepository == nil {
		s.inventoryRepository = memory.NewInventoryRepository()
	}

	return s.inventoryRepository
}

// ServerImplementation creates GRPC service handler
// Initializes all required dependencies (repository)
func (s *serviceProvider) ServerImplementation(ctx context.Context) *server.Implementation {
	if s.serverImplementation == nil {
		inventoryRepository := s.InventoryRepository(ctx)
		s.serverImplementation = server.NewImplementation(inventoryRepository)
	}

	return s.serverImplementation
}
