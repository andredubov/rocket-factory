package memory

import (
	"context"
	"fmt"

	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// UpdatePart modifies an existing part in the repository
// Thread-safe operation using mutex lock
// Returns:
// - nil if part was updated successfully
// - error if part with specified UUID doesn't exist
func (p *inventoryRepository) UpdatePart(ctx context.Context, part model.Part) error {
	p.mu.Lock()         // Acquire write lock
	defer p.mu.Unlock() // Ensure lock is released

	// Verify part exists before update
	if _, exists := p.parts[part.Uuid]; !exists {
		return fmt.Errorf("part with UUID %s not found", part.Uuid)
	}

	// Create defensive copy to prevent external modifications
	updatedPart := part
	p.parts[part.Uuid] = &updatedPart
	return nil
}
