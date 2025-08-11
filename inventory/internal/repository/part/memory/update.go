package memory

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
)

// UpdatePart modifies an existing part in the repository
func (r *inventoryRepository) UpdatePart(ctx context.Context, part model.Part) error {
	r.mu.Lock()         // Acquire write lock
	defer r.mu.Unlock() // Ensure lock is released

	repoPart := converter.PartToRepoModel(part)
	// Verify part exists before update
	if _, exists := r.parts[repoPart.Uuid]; !exists {
		return model.ErrPartWithUUIDNotFound(repoPart.Uuid)
	}

	// Create defensive copy to prevent external modifications
	updatedPart := repoPart
	r.parts[part.Uuid] = &updatedPart
	return nil
}
