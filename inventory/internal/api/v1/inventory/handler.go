package server

import (
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// InventoryImplementation is the gRPC server implementation for the InventoryService.
type InventoryImplementation struct {
	inventory_v1.InventoryServiceServer                      // Embedded gRPC service interface
	inventoryRepository                 repository.Inventory // Repository for data access
}

// NewInventoryImplementation creates a new instance of the gRPC server implementation.
func NewInventoryImplementation(repository repository.Inventory) *InventoryImplementation {
	return &InventoryImplementation{
		inventoryRepository: repository,
	}
}
