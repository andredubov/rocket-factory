package memory

import (
	"strings"
	"sync"

	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// inventoryRepository is an in-memory implementation of the Inventory repository
// Uses a map for storage and RWMutex for concurrent access safety
type inventoryRepository struct {
	mu    sync.RWMutex           // Read-write mutex to protect concurrent access
	parts map[string]*model.Part // Map storing parts by their UUID
}

// NewInventoryRepository creates a new in-memory inventory repository instance
// Returns a ready-to-use repository with initialized storage map
func NewInventoryRepository() repository.Inventory {
	return &inventoryRepository{
		parts: make(map[string]*model.Part), // Initialize empty parts map
	}
}

// Filter by tags (OR logic within tags, part must have at least one matching tag)
func filterByTags(parts []model.Part, tags []string) []model.Part {
	var result []model.Part
	for _, part := range parts {
		partMatched := false
		for _, tag := range tags {
			for _, partTag := range part.Tags {
				if strings.EqualFold(partTag, tag) {
					result = append(result, part)
					partMatched = true
					break
				}
			}
			if partMatched {
				break
			}
		}
	}
	return result
}

// Filter by manufacturer country (OR logic within countries)
func filterByCountry(parts []model.Part, countries []string) []model.Part {
	var result []model.Part
	for _, part := range parts {
		for _, country := range countries {
			if strings.EqualFold(part.Manufacturer.Country, country) {
				result = append(result, part)
				break
			}
		}
	}
	return result
}

// Filter by category (OR logic within categories)
func filterByCategory(parts []model.Part, categories []model.Category) []model.Part {
	var result []model.Part
	for _, part := range parts {
		for _, cat := range categories {
			if part.Category.ID == cat.ID || strings.EqualFold(part.Category.Name, cat.Name) {
				result = append(result, part)
				break
			}
		}
	}
	return result
}

// Filter by name (OR logic within names)
func filterByName(parts []model.Part, names []string) []model.Part {
	var result []model.Part
	for _, part := range parts {
		for _, name := range names {
			if strings.EqualFold(part.Name, name) {
				result = append(result, part)
				break
			}
		}
	}
	return result
}

// Helper function to check if filter is empty
func isEmptyFilter(filter model.PartFilter) bool {
	return len(filter.UUIDs) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}
