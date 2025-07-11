package memory

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/repository"
)

// DeletePart removes a part from the repository by UUID
// Thread-safe operation using mutex lock
// Returns:
// - nil if part was deleted successfully
// - error if part with specified UUID doesn't exist
func (p *inventoryRepository) DeletePart(ctx context.Context, uuid string) error {
	p.mu.Lock()         // Acquire write lock
	defer p.mu.Unlock() // Ensure lock is released

	// Verify part exists before deletion
	if _, exists := p.parts[uuid]; !exists {
		return repository.ErrPartWithUUIDNotFound(uuid)
	}

	delete(p.parts, uuid) // Remove part from map
	return nil
}
