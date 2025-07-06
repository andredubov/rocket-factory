package server

import (
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// Implementation is the gRPC server implementation for the InventoryService.
// It embeds the generated InventoryServiceServer interface and contains
// the business logic for handling inventory operations.
type Implementation struct {
	inventory_v1.InventoryServiceServer                      // Embedded gRPC service interface
	inventoryRepository                 repository.Inventory // Repository for data access
}

// NewImplementation creates a new instance of the gRPC server implementation.
// It requires an Inventory repository interface to be passed as a dependency
// for data persistence operations.
func NewImplementation(repository repository.Inventory) *Implementation {
	return &Implementation{
		inventoryRepository: repository,
	}
}
