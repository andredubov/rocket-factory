package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// Error definitions for the inventory repository.
// These are base errors that can be wrapped with additional context.
var (
	ErrPartAlreadyExists = errors.New("part already exists")
	ErrPartNotFound      = errors.New("part not found")
)

// ErrPartWithUUIDNotFound constructs a new error indicating a part with the specified UUID was not found.
func ErrPartWithUUIDNotFound(uuid string) error {
	return fmt.Errorf("part with UUID %s not found: %w", uuid, ErrPartNotFound)
}

// ErrPartWithUUIDExists constructs a new error indicating a part with the specified UUID already exists.
func ErrPartWithUUIDExists(uuid string) error {
	return fmt.Errorf("part with UUID %s already exists: %w", uuid, ErrPartAlreadyExists)
}

// Inventory defines the interface for inventory repository operations.
// All implementations must provide thread-safe access to the underlying data store.
type Inventory interface {
	GetPartList(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	AddPart(ctx context.Context, part model.Part) error
	UpdatePart(ctx context.Context, part model.Part) error
	DeletePart(ctx context.Context, uuid string) error
}
