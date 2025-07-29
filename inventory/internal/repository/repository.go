package repository

import (
	"fmt"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
)

// ErrPartWithUUIDNotFound constructs a new error indicating a part with the specified UUID was not found.
func ErrPartWithUUIDNotFound(uuid string) error {
	return fmt.Errorf("part with UUID %s not found: %w", uuid, model.ErrPartNotFound)
}

// ErrPartWithUUIDExists constructs a new error indicating a part with the specified UUID already exists.
func ErrPartWithUUIDExists(uuid string) error {
	return fmt.Errorf("part with UUID %s already exists: %w", uuid, model.ErrPartAlreadyExists)
}
