package converter

import (
	"github.com/andredubov/rocket-factory/inventory/internal/model"
	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// PartFilterToRepoModel converts a domain model (repo.PartFilter) to a repository model (repoModel.PartFilter).
func PartFilterToRepoModel(partFilter model.PartFilter) repoModel.PartFilter {
	// Convert categories from model to repoModel
	var repoCategories []repoModel.PartCategory
	for _, category := range partFilter.Categories {
		repoCategories = append(repoCategories, repoModel.PartCategory(category))
	}

	return repoModel.PartFilter{
		UUIDs:                 partFilter.UUIDs,
		Names:                 partFilter.Names,
		Categories:            repoCategories,
		ManufacturerCountries: partFilter.ManufacturerCountries,
		Tags:                  partFilter.Tags,
	}
}
