package server

import (
	"github.com/andredubov/rocket-factory/inventory/internal/service"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// InventoryImplementation is the gRPC server implementation for the InventoryService.
type InventoryImplementation struct {
	inventory_v1.InventoryServiceServer                   // Embedded gRPC service interface
	inventoryService                    service.Inventory // Inventory service
}

// NewInventoryImplementation creates a new instance of the gRPC server implementation.
func NewInventoryImplementation(service service.Inventory) *InventoryImplementation {
	return &InventoryImplementation{
		inventoryService: service,
	}
}
