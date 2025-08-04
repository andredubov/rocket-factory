package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// TestPartFilterToRepoModel_CompleteConversion verifies the complete conversion of a PartFilter
// from domain model to repository model with all filter fields populated. Tests that all filter
// criteria including UUIDs, names, categories, countries and tags are correctly mapped.
func TestPartFilterToRepoModel_CompleteConversion(t *testing.T) {
	// Setup
	source := model.PartFilter{
		UUIDs:                 []string{gofakeit.UUID(), gofakeit.UUID()},
		Names:                 []string{gofakeit.Word(), gofakeit.Word()},
		Categories:            []model.PartCategory{model.PartCategoryEngine, model.PartCategoryFuel},
		ManufacturerCountries: []string{gofakeit.Country(), gofakeit.Country()},
		Tags:                  []string{gofakeit.Word(), gofakeit.Word()},
	}

	// Test
	result := converter.PartFilterToRepoModel(source)

	// Verify
	require.Equal(t, source.UUIDs, result.UUIDs)
	require.Equal(t, source.Names, result.Names)
	require.Equal(t, len(source.Categories), len(result.Categories))
	for i, category := range source.Categories {
		require.Equal(t, repoModel.PartCategory(category), result.Categories[i])
	}
	require.Equal(t, source.ManufacturerCountries, result.ManufacturerCountries)
	require.Equal(t, source.Tags, result.Tags)
}

// TestPartFilterToRepoModel_EmptyFilter verifies the conversion handles an empty filter correctly.
// Tests that when all filter fields are empty/unset, the converted repository model maintains
// nil slices rather than empty slices for all filter criteria.
func TestPartFilterToRepoModel_EmptyFilter(t *testing.T) {
	// Setup
	source := model.PartFilter{}

	// Test
	result := converter.PartFilterToRepoModel(source)

	// Verify
	require.Nil(t, result.UUIDs)
	require.Nil(t, result.Names)
	require.Nil(t, result.Categories)
	require.Nil(t, result.ManufacturerCountries)
	require.Nil(t, result.Tags)
}

// TestPartFilterToRepoModel_PartialFields verifies the conversion works correctly with only
// some filter fields set. Tests that unset fields remain nil while set fields are properly
// converted, including proper category enum mapping for the set fields.
func TestPartFilterToRepoModel_PartialFields(t *testing.T) {
	// Setup
	source := model.PartFilter{
		Categories: []model.PartCategory{model.PartCategoryPorthole},
		Tags:       []string{"special-tag"},
	}

	// Test
	result := converter.PartFilterToRepoModel(source)

	// Verify
	require.Nil(t, result.UUIDs)
	require.Nil(t, result.Names)
	require.Equal(t, 1, len(result.Categories))
	require.Equal(t, repoModel.PartCategoryPorthole, result.Categories[0])
	require.Nil(t, result.ManufacturerCountries)
	require.Equal(t, []string{"special-tag"}, result.Tags)
}

// TestPartFilterToRepoModel_AllCategories specifically tests category filter conversion
// with all possible category values. Verifies the complete set of part categories are
// correctly mapped between domain and repository models.
func TestPartFilterToRepoModel_AllCategories(t *testing.T) {
	// Setup
	source := model.PartFilter{
		Categories: []model.PartCategory{
			model.PartCategoryUnknown,
			model.PartCategoryEngine,
			model.PartCategoryFuel,
			model.PartCategoryPorthole,
			model.PartCategoryWing,
		},
	}

	// Test
	result := converter.PartFilterToRepoModel(source)

	// Verify
	require.Equal(t, 5, len(result.Categories))
	require.Equal(t, repoModel.PartCategoryUnknown, result.Categories[0])
	require.Equal(t, repoModel.PartCategoryEngine, result.Categories[1])
	require.Equal(t, repoModel.PartCategoryFuel, result.Categories[2])
	require.Equal(t, repoModel.PartCategoryPorthole, result.Categories[3])
	require.Equal(t, repoModel.PartCategoryWing, result.Categories[4])
}
