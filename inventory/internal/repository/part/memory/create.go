package memory

import (
	"context"
	"fmt"

	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// AddPart adds a new part to the in-memory repository
// Thread-safe operation using mutex lock
// Returns:
// - nil if part was added successfully
// - error if part with same UUID already exists
func (p *inventoryRepository) AddPart(ctx context.Context, part model.Part) error {
	p.mu.Lock()         // Acquire write lock
	defer p.mu.Unlock() // Ensure lock is released

	// Check for existing part with same UUID
	if _, exists := p.parts[part.Uuid]; exists {
		return fmt.Errorf("part with UUID %s already exists", part.Uuid)
	}

	// Create defensive copy to prevent external modifications
	newPart := part
	p.parts[part.Uuid] = &newPart
	return nil
}
