package memory

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
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
func (r *inventoryRepository) GetPartList(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	r.mu.RLock()         // Acquire read lock
	defer r.mu.RUnlock() // Ensure lock is released

	repoFilter := converter.PartFilterToRepoModel(filter)

	// Return all parts if no filters specified
	if isEmptyFilter(repoFilter) {
		parts := make([]model.Part, 0, len(r.parts))
		for _, repoPart := range r.parts {
			part := converter.PartToModel(*repoPart)
			parts = append(parts, part)
		}
		return parts, nil
	}

	var repoPartList []repoModel.Part

	// First filter pass - by UUIDs (OR logic)
	if len(repoFilter.UUIDs) > 0 {
		for _, uuid := range repoFilter.UUIDs {
			if repoPart, exists := r.parts[uuid]; exists {
				repoPartList = append(repoPartList, *repoPart)
			}
		}
	} else {
		// If no UUID filter, start with all parts
		for _, repoPart := range r.parts {
			repoPartList = append(repoPartList, *repoPart)
		}
	}

	// Apply subsequent filters (AND logic between fields)
	if len(filter.Names) > 0 {
		repoPartList = filterByName(repoPartList, repoFilter.Names)
	}
	if len(filter.Categories) > 0 {
		repoPartList = filterByCategory(repoPartList, repoFilter.Categories)
	}
	if len(filter.ManufacturerCountries) > 0 {
		repoPartList = filterByCountry(repoPartList, repoFilter.ManufacturerCountries)
	}
	if len(filter.Tags) > 0 {
		repoPartList = filterByTags(repoPartList, repoFilter.Tags)
	}

	partList := make([]model.Part, 0, len(repoPartList))
	for _, repoPart := range repoPartList {
		partList = append(partList, converter.PartToModel(repoPart))
	}

	return partList, nil
}
