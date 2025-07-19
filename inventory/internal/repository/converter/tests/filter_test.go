package tests

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// TestPartFilterToRepoModel_CompleteConversion verifies the complete conversion of a PartFilter
// from domain model to repository model with all filter fields populated. Tests that all filter
// criteria including UUIDs, names, categories, countries and tags are correctly mapped.
func (s *InventoryRepoConverterSuite) TestPartFilterToRepoModel_CompleteConversion() {
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
	s.Equal(source.UUIDs, result.UUIDs)
	s.Equal(source.Names, result.Names)
	s.Equal(len(source.Categories), len(result.Categories))
	for i, category := range source.Categories {
		s.Equal(repoModel.PartCategory(category), result.Categories[i])
	}
	s.Equal(source.ManufacturerCountries, result.ManufacturerCountries)
	s.Equal(source.Tags, result.Tags)
}

// TestPartFilterToRepoModel_EmptyFilter verifies the conversion handles an empty filter correctly.
// Tests that when all filter fields are empty/unset, the converted repository model maintains
// nil slices rather than empty slices for all filter criteria.
func (s *InventoryRepoConverterSuite) TestPartFilterToRepoModel_EmptyFilter() {
	// Setup
	source := model.PartFilter{}

	// Test
	result := converter.PartFilterToRepoModel(source)

	// Verify
	s.Require().Nil(result.UUIDs)
	s.Require().Nil(result.Names)
	s.Require().Nil(result.Categories)
	s.Require().Nil(result.ManufacturerCountries)
	s.Require().Nil(result.Tags)
}

// TestPartFilterToRepoModel_PartialFields verifies the conversion works correctly with only
// some filter fields set. Tests that unset fields remain nil while set fields are properly
// converted, including proper category enum mapping for the set fields.
func (s *InventoryRepoConverterSuite) TestPartFilterToRepoModel_PartialFields() {
	// Setup
	source := model.PartFilter{
		Categories: []model.PartCategory{model.PartCategoryPorthole},
		Tags:       []string{"special-tag"},
	}

	// Test
	result := converter.PartFilterToRepoModel(source)

	// Verify
	s.Require().Nil(result.UUIDs)
	s.Require().Nil(result.Names)
	s.Require().Equal(1, len(result.Categories))
	s.Require().Equal(repoModel.PartCategoryPorthole, result.Categories[0])
	s.Require().Nil(result.ManufacturerCountries)
	s.Require().Equal([]string{"special-tag"}, result.Tags)
}

// TestPartFilterToRepoModel_AllCategories specifically tests category filter conversion
// with all possible category values. Verifies the complete set of part categories are
// correctly mapped between domain and repository models.
func (s *InventoryRepoConverterSuite) TestPartFilterToRepoModel_AllCategories() {
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
	s.Require().Equal(5, len(result.Categories))
	s.Require().Equal(repoModel.PartCategoryUnknown, result.Categories[0])
	s.Require().Equal(repoModel.PartCategoryEngine, result.Categories[1])
	s.Require().Equal(repoModel.PartCategoryFuel, result.Categories[2])
	s.Require().Equal(repoModel.PartCategoryPorthole, result.Categories[3])
	s.Require().Equal(repoModel.PartCategoryWing, result.Categories[4])
}
