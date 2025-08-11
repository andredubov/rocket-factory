package memory

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
)

// DeletePart removes a part from the repository by UUID
// Thread-safe operation using mutex lock
// Returns:
// - nil if part was deleted successfully
// - error if part with specified UUID doesn't exist
func (r *inventoryRepository) DeletePart(ctx context.Context, uuid string) error {
	r.mu.Lock()         // Acquire write lock
	defer r.mu.Unlock() // Ensure lock is released

	// Verify part exists before deletion
	if _, exists := r.parts[uuid]; !exists {
		return model.ErrPartWithUUIDNotFound(uuid)
	}

	delete(r.parts, uuid) // Remove part from map
	return nil
}
