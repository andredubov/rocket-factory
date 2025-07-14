package repository

import (
	"context"
	"fmt"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// ErrPartWithUUIDNotFound constructs a new error indicating a part with the specified UUID was not found.
func ErrPartWithUUIDNotFound(uuid string) error {
	return fmt.Errorf("part with UUID %s not found: %w", uuid, model.ErrPartNotFound)
}

// ErrPartWithUUIDExists constructs a new error indicating a part with the specified UUID already exists.
func ErrPartWithUUIDExists(uuid string) error {
	return fmt.Errorf("part with UUID %s already exists: %w", uuid, model.ErrPartAlreadyExists)
}

// Inventory defines the interface for inventory repository operations.
type Inventory interface {
	GetPartList(ctx context.Context, filter repoModel.PartFilter) ([]repoModel.Part, error)
	GetPart(ctx context.Context, uuid string) (*repoModel.Part, error)
	AddPart(ctx context.Context, part repoModel.Part) error
	UpdatePart(ctx context.Context, part repoModel.Part) error
	DeletePart(ctx context.Context, uuid string) error
}
