package memory

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
)

// UpdatePart modifies an existing part in the repository
func (p *inventoryRepository) UpdatePart(ctx context.Context, part model.Part) error {
	p.mu.Lock()         // Acquire write lock
	defer p.mu.Unlock() // Ensure lock is released

	repoPart := converter.PartToRepoModel(part)
	// Verify part exists before update
	if _, exists := p.parts[repoPart.Uuid]; !exists {
		return repository.ErrPartWithUUIDNotFound(repoPart.Uuid)
	}

	// Create defensive copy to prevent external modifications
	updatedPart := repoPart
	p.parts[part.Uuid] = &updatedPart
	return nil
}
