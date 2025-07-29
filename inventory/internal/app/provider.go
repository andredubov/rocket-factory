package app

import (
	"context"
	"log"

	"github.com/andredubov/golibs/pkg/config"
	"github.com/andredubov/golibs/pkg/config/env"

	api "github.com/andredubov/rocket-factory/inventory/internal/api/v1/inventory"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/memory"
	"github.com/andredubov/rocket-factory/inventory/internal/service"
	"github.com/andredubov/rocket-factory/inventory/internal/service/inventory"
)

// serviceProvider implements the dependency container pattern
// It provides lazy initialization of application components
type serviceProvider struct {
	inventoryRepository  service.InventoryRepository
	inventoryService     api.InventoryService
	grpcConfig           config.GRPCConfig            // GRPC server configuration
	serverImplementation *api.InventoryImplementation // GRPC service implementation
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
func (s *serviceProvider) InventoryRepository(ctx context.Context) service.InventoryRepository {
	if s.inventoryRepository == nil {
		s.inventoryRepository = memory.NewInventoryRepository()
	}

	return s.inventoryRepository
}

// InventoryService provides access to inventory service layer
func (s *serviceProvider) InventoryService(ctx context.Context) api.InventoryService {
	if s.inventoryService == nil {
		s.inventoryService = inventory.NewService(
			s.InventoryRepository(ctx),
		)
	}

	return s.inventoryService
}

// ServerImplementation creates GRPC service handler
// Initializes all required dependencies (service)
func (s *serviceProvider) ServerImplementation(ctx context.Context) *api.InventoryImplementation {
	if s.serverImplementation == nil {
		inventoryService := s.InventoryService(ctx)
		s.serverImplementation = api.NewInventoryImplementation(inventoryService)
	}

	return s.serverImplementation
}
