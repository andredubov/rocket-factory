package memory

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
)

// GetPart retrieves a single part by UUID
// Thread-safe read operation using RWMutex
// Returns:
// - Part if found
// - error if part doesn't exist
func (r *inventoryRepository) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()         // Acquire read lock
	defer r.mu.RUnlock() // Ensure lock is released

	repoPart, exists := r.parts[uuid]
	if !exists {
		return nil, model.ErrPartWithUUIDNotFound(uuid)
	}

	part := converter.PartToModel(*repoPart)
	return &part, nil
}
