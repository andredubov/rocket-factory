package memory

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// GetPartList retrieves parts matching the filter criteria
// Thread-safe read operation using RWMutex
// Filtering logic:
// - Empty filter returns all parts
// - OR logic within each filter field
// - AND logic between different filter fields
// Returns:
// - Slice of matching parts
// - nil error if successful
func (i *inventoryRepository) GetPartList(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	i.mu.RLock()         // Acquire read lock
	defer i.mu.RUnlock() // Ensure lock is released

	// Return all parts if no filters specified
	if isEmptyFilter(filter) {
		parts := make([]model.Part, 0, len(i.parts))
		for _, part := range i.parts {
			parts = append(parts, *part)
		}
		return parts, nil
	}

	var result []model.Part

	// First filter pass - by UUIDs (OR logic)
	if len(filter.UUIDs) > 0 {
		for _, uuid := range filter.UUIDs {
			if part, exists := i.parts[uuid]; exists {
				result = append(result, *part)
			}
		}
	} else {
		// If no UUID filter, start with all parts
		for _, part := range i.parts {
			result = append(result, *part)
		}
	}

	// Apply subsequent filters (AND logic between fields)
	if len(filter.Names) > 0 {
		result = filterByName(result, filter.Names)
	}
	if len(filter.Categories) > 0 {
		result = filterByCategory(result, filter.Categories)
	}
	if len(filter.ManufacturerCountries) > 0 {
		result = filterByCountry(result, filter.ManufacturerCountries)
	}
	if len(filter.Tags) > 0 {
		result = filterByTags(result, filter.Tags)
	}

	return result, nil
}

// GetPart retrieves a single part by UUID
// Thread-safe read operation using RWMutex
// Returns:
// - Part if found
// - error if part doesn't exist
func (i *inventoryRepository) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	i.mu.RLock()         // Acquire read lock
	defer i.mu.RUnlock() // Ensure lock is released

	part, exists := i.parts[uuid]
	if !exists {
		return nil, repository.ErrPartWithUUIDNotFound(uuid)
	}

	return part, nil
}
