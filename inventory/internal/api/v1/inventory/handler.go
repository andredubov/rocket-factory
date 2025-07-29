package api

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// InventoryService defines the interface for inventory service operations.
type InventoryService interface {
	GetPartList(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	AddPart(ctx context.Context, part model.Part) error
	UpdatePart(ctx context.Context, part model.Part) error
	DeletePart(ctx context.Context, uuid string) error
}

// InventoryImplementation is the gRPC server implementation for the InventoryService.
type InventoryImplementation struct {
	inventory_v1.InventoryServiceServer                  // Embedded gRPC service interface
	inventoryService                    InventoryService // Inventory service
}

// NewInventoryImplementation creates a new instance of the gRPC server implementation.
func NewInventoryImplementation(service InventoryService) *InventoryImplementation {
	return &InventoryImplementation{
		inventoryService: service,
	}
}
