package service

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
)

// InventoryRepository defines the interface for inventory repository operations.
type InventoryRepository interface {
	GetPartList(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	AddPart(ctx context.Context, part model.Part) error
	UpdatePart(ctx context.Context, part model.Part) error
	DeletePart(ctx context.Context, uuid string) error
}
