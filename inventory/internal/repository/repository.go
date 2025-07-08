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
	// ErrPartAlreadyExists is returned when attempting to add a part with a UUID that already exists.
	ErrPartAlreadyExists = errors.New("part already exists")

	// ErrPartNotFound is returned when a requested part cannot be found.
	ErrPartNotFound = errors.New("part not found")
)

// ErrPartWithUUIDNotFound constructs a new error indicating a part with the specified UUID was not found.
// The error wraps the base ErrPartNotFound for error handling and provides the specific UUID in the message.
//
// Parameters:
//   - uuid: The UUID of the part that was not found
//
// Returns:
//   - error: Formatted error with UUID context, wrapping ErrPartNotFound
func ErrPartWithUUIDNotFound(uuid string) error {
	return fmt.Errorf("part with UUID %s not found: %w", uuid, ErrPartNotFound)
}

// ErrPartWithUUIDExists constructs a new error indicating a part with the specified UUID already exists.
// The error wraps the base ErrPartAlreadyExists for error handling and provides the specific UUID in the message.
//
// Parameters:
//   - uuid: The UUID of the part that already exists
//
// Returns:
//   - error: Formatted error with UUID context, wrapping ErrPartAlreadyExists
func ErrPartWithUUIDExists(uuid string) error {
	return fmt.Errorf("part with UUID %s already exists: %w", uuid, ErrPartAlreadyExists)
}

// Inventory defines the interface for inventory repository operations.
// All implementations must provide thread-safe access to the underlying data store.
type Inventory interface {
	// GetPartList retrieves a list of parts matching the filter criteria.
	// Filter fields are combined with AND logic, while values within each field use OR logic.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts
	//   - filter: Criteria for filtering parts
	//
	// Returns:
	//   - []model.Part: Slice of matching parts (empty slice if no matches)
	//   - error: Repository error if operation fails
	GetPartList(ctx context.Context, filter model.PartFilter) ([]model.Part, error)

	// GetPart retrieves a single part by its UUID.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts
	//   - uuid: Unique identifier of the part
	//
	// Returns:
	//   - *model.Part: Pointer to the found part
	//   - error: ErrPartNotFound if part doesn't exist
	GetPart(ctx context.Context, uuid string) (*model.Part, error)

	// AddPart inserts a new part into the repository.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts
	//   - part: Part to be added
	//
	// Returns:
	//   - error: ErrPartAlreadyExists if part with UUID exists, or nil if part was added successfully
	AddPart(ctx context.Context, part model.Part) error

	// UpdatePart modifies an existing part in the repository.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts
	//   - part: Part with updated values (identified by UUID)
	//
	// Returns:
	//   - error: ErrPartNotFound if part doesn't exist, or nil if part was updated successfully
	UpdatePart(ctx context.Context, part model.Part) error

	// DeletePart removes a part from the repository by UUID.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeouts
	//   - uuid: Unique identifier of the part to delete
	//
	// Returns:
	//   - error: ErrPartNotFound if part doesn't exist, or nil if part was deleted successfully
	DeletePart(ctx context.Context, uuid string) error
}
