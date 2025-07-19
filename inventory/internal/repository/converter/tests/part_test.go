package tests

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// Helper functions to create pointers
func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func boolPtr(b bool) *bool {
	return &b
}

// TestPartToRepoModel_CompleteConversion verifies the complete conversion from domain model
// to repository model with all fields populated. Tests that all field values including nested
// structures, metadata, and timestamps are correctly mapped.
func (s *InventoryRepoConverterSuite) TestPartToRepoModel_CompleteConversion() {
	// Setup
	now := time.Now()
	source := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Word(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Float64Range(1, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      model.PartCategory(gofakeit.IntRange(1, 4)),
		Dimensions: model.Dimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(1, 100),
		},
		Manufacturer: model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: []string{gofakeit.Word(), gofakeit.Word()},
		Metadata: map[string]model.Value{
			"string":  {StringValue: stringPtr("value")},
			"int":     {Int64Value: int64Ptr(42)},
			"float":   {DoubleValue: float64Ptr(3.14)},
			"boolean": {BoolValue: boolPtr(true)},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test
	result := converter.PartToRepoModel(source)

	// Verify
	s.Require().Equal(source.Uuid, result.Uuid)
	s.Require().Equal(source.Name, result.Name)
	s.Require().Equal(source.Description, result.Description)
	s.Require().Equal(source.Price, result.Price)
	s.Require().Equal(source.StockQuantity, result.StockQuantity)
	s.Require().Equal(repoModel.PartCategory(source.Category), result.Category)

	// Dimensions
	s.Require().Equal(source.Dimensions.Length, result.Dimensions.Length)
	s.Require().Equal(source.Dimensions.Width, result.Dimensions.Width)
	s.Require().Equal(source.Dimensions.Height, result.Dimensions.Height)
	s.Require().Equal(source.Dimensions.Weight, result.Dimensions.Weight)

	// Manufacturer
	s.Require().Equal(source.Manufacturer.Name, result.Manufacturer.Name)
	s.Require().Equal(source.Manufacturer.Country, result.Manufacturer.Country)
	s.Require().Equal(source.Manufacturer.Website, result.Manufacturer.Website)

	// Tags
	s.Require().Equal(source.Tags, result.Tags)

	// Metadata
	s.Require().NotNil(result.Metadata)
	s.Require().Equal("value", *result.Metadata["string"].StringValue)
	s.Require().Equal(int64(42), *result.Metadata["int"].Int64Value)
	s.Require().Equal(3.14, *result.Metadata["float"].DoubleValue)
	s.Require().Equal(true, *result.Metadata["boolean"].BoolValue)

	// Timestamps
	s.Require().Equal(source.CreatedAt, result.CreatedAt)
	s.Require().Equal(source.UpdatedAt, result.UpdatedAt)
}

// TestPartToRepoModel_EmptyFields verifies the conversion handles empty/zero values correctly.
// Tests that empty strings, zero values, nil slices, and empty nested structures are properly
// converted without errors.
func (s *InventoryRepoConverterSuite) TestPartToRepoModel_EmptyFields() {
	// Setup
	source := model.Part{
		Uuid:          "",
		Name:          "",
		Description:   "",
		Price:         0,
		StockQuantity: 0,
		Category:      model.PartCategoryUnknown,
		Dimensions:    model.Dimensions{},
		Manufacturer:  model.Manufacturer{},
		Tags:          nil,
		Metadata:      nil,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}

	// Test
	result := converter.PartToRepoModel(source)

	// Verify
	s.Require().Equal(source.Uuid, result.Uuid)
	s.Require().Equal(source.Name, result.Name)
	s.Require().Equal(source.Description, result.Description)
	s.Require().Equal(source.Price, result.Price)
	s.Require().Equal(source.StockQuantity, result.StockQuantity)
	s.Require().Equal(repoModel.PartCategory(source.Category), result.Category)
	s.Require().Equal(source.Dimensions.Length, result.Dimensions.Length)
	s.Require().Equal(source.Dimensions.Width, result.Dimensions.Width)
	s.Require().Equal(source.Dimensions.Height, result.Dimensions.Height)
	s.Require().Equal(source.Dimensions.Weight, result.Dimensions.Weight)
	s.Require().Equal(source.Manufacturer.Name, result.Manufacturer.Name)
	s.Require().Equal(source.Manufacturer.Country, result.Manufacturer.Country)
	s.Require().Equal(source.Manufacturer.Website, result.Manufacturer.Website)
	s.Require().Nil(result.Tags)
	s.Require().NotNil(result.Metadata)
	s.Require().Equal(source.CreatedAt, result.CreatedAt)
	s.Require().Equal(source.UpdatedAt, result.UpdatedAt)
}

// TestPartToRepoModel_PartialMetadata verifies metadata conversion handles partial/null values.
// Tests that metadata fields with some nil values are correctly converted while maintaining
// the structure of non-nil values.
func (s *InventoryRepoConverterSuite) TestPartToRepoModel_PartialMetadata() {
	// Setup
	source := model.Part{
		Metadata: map[string]model.Value{
			"string": {StringValue: stringPtr("test")},
			"int":    {Int64Value: nil},
		},
	}

	// Test
	result := converter.PartToRepoModel(source)

	// Verify
	s.Require().NotNil(result.Metadata)
	s.Require().Equal("test", *result.Metadata["string"].StringValue)
	s.Require().Nil(result.Metadata["string"].Int64Value)
	s.Require().Nil(result.Metadata["string"].DoubleValue)
	s.Require().Nil(result.Metadata["string"].BoolValue)
	s.Require().Nil(result.Metadata["int"].StringValue)
	s.Require().Nil(result.Metadata["int"].Int64Value)
}

// TestPartToModel_CompleteConversion verifies the complete conversion from repository model
// to domain model with all fields populated. Tests the reverse mapping of all field values
// including nested structures and metadata.
func (s *InventoryRepoConverterSuite) TestPartToModel_CompleteConversion() {
	// Setup
	now := time.Now()
	source := repoModel.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Word(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Float64Range(1, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      repoModel.PartCategory(gofakeit.IntRange(1, 4)),
		Dimensions: repoModel.Dimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(1, 100),
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: []string{gofakeit.Word(), gofakeit.Word()},
		Metadata: map[string]repoModel.Value{
			"string":  {StringValue: stringPtr("value")},
			"int":     {Int64Value: int64Ptr(42)},
			"float":   {DoubleValue: float64Ptr(3.14)},
			"boolean": {BoolValue: boolPtr(true)},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test
	result := converter.PartToModel(source)

	// Verify
	s.Require().Equal(source.Uuid, result.Uuid)
	s.Require().Equal(source.Name, result.Name)
	s.Require().Equal(source.Description, result.Description)
	s.Require().Equal(source.Price, result.Price)
	s.Require().Equal(source.StockQuantity, result.StockQuantity)
	s.Require().Equal(model.PartCategory(source.Category), result.Category)

	// Dimensions
	s.Require().Equal(source.Dimensions.Length, result.Dimensions.Length)
	s.Require().Equal(source.Dimensions.Width, result.Dimensions.Width)
	s.Require().Equal(source.Dimensions.Height, result.Dimensions.Height)
	s.Require().Equal(source.Dimensions.Weight, result.Dimensions.Weight)

	// Manufacturer
	s.Require().Equal(source.Manufacturer.Name, result.Manufacturer.Name)
	s.Require().Equal(source.Manufacturer.Country, result.Manufacturer.Country)
	s.Require().Equal(source.Manufacturer.Website, result.Manufacturer.Website)

	// Tags
	s.Require().Equal(source.Tags, result.Tags)

	// Metadata
	s.Require().NotNil(result.Metadata)
	s.Require().Equal("value", *result.Metadata["string"].StringValue)
	s.Require().Equal(int64(42), *result.Metadata["int"].Int64Value)
	s.Require().Equal(3.14, *result.Metadata["float"].DoubleValue)
	s.Require().Equal(true, *result.Metadata["boolean"].BoolValue)

	// Timestamps
	s.Require().Equal(source.CreatedAt, result.CreatedAt)
	s.Require().Equal(source.UpdatedAt, result.UpdatedAt)
}

// TestPartToModel_EmptyFields verifies the reverse conversion handles empty/zero values correctly.
// Ensures empty repository model fields are properly converted to their domain model equivalents.
func (s *InventoryRepoConverterSuite) TestPartToModel_EmptyFields() {
	// Setup
	source := repoModel.Part{
		Uuid:          "",
		Name:          "",
		Description:   "",
		Price:         0,
		StockQuantity: 0,
		Category:      repoModel.PartCategoryUnknown,
		Dimensions:    repoModel.Dimensions{},
		Manufacturer:  repoModel.Manufacturer{},
		Tags:          nil,
		Metadata:      nil,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}

	// Test
	result := converter.PartToModel(source)

	// Verify
	s.Require().Equal(source.Uuid, result.Uuid)
	s.Require().Equal(source.Name, result.Name)
	s.Require().Equal(source.Description, result.Description)
	s.Require().Equal(source.Price, result.Price)
	s.Require().Equal(source.StockQuantity, result.StockQuantity)
	s.Require().Equal(model.PartCategory(source.Category), result.Category)
	s.Require().Equal(source.Dimensions.Length, result.Dimensions.Length)
	s.Require().Equal(source.Dimensions.Width, result.Dimensions.Width)
	s.Require().Equal(source.Dimensions.Height, result.Dimensions.Height)
	s.Require().Equal(source.Dimensions.Weight, result.Dimensions.Weight)
	s.Require().Equal(source.Manufacturer.Name, result.Manufacturer.Name)
	s.Require().Equal(source.Manufacturer.Country, result.Manufacturer.Country)
	s.Require().Equal(source.Manufacturer.Website, result.Manufacturer.Website)
	s.Require().Nil(result.Tags)
	s.Require().NotNil(result.Metadata)
	s.Require().Equal(source.CreatedAt, result.CreatedAt)
	s.Require().Equal(source.UpdatedAt, result.UpdatedAt)
}

// TestPartToModel_PartialMetadata verifies reverse metadata conversion handles partial values.
// Tests that repository metadata with some nil values converts correctly to domain model metadata.
func (s *InventoryRepoConverterSuite) TestPartToModel_PartialMetadata() {
	// Setup
	source := repoModel.Part{
		Metadata: map[string]repoModel.Value{
			"string": {StringValue: stringPtr("test")},
			"int":    {Int64Value: nil},
		},
	}

	// Test
	result := converter.PartToModel(source)

	// Verify
	s.Require().NotNil(result.Metadata)
	s.Require().Equal("test", *result.Metadata["string"].StringValue)
	s.Require().Nil(result.Metadata["string"].Int64Value)
	s.Require().Nil(result.Metadata["string"].DoubleValue)
	s.Require().Nil(result.Metadata["string"].BoolValue)
	s.Require().Nil(result.Metadata["int"].StringValue)
	s.Require().Nil(result.Metadata["int"].Int64Value)
}

// TestPartToModel_CategoryConversion specifically tests the category enum conversion between
// repository and domain models. Verifies all valid category values are correctly mapped.
func (s *InventoryRepoConverterSuite) TestPartToModel_CategoryConversion() {
	// Test all valid category values
	categories := []repoModel.PartCategory{
		repoModel.PartCategoryUnknown,
		repoModel.PartCategoryEngine,
		repoModel.PartCategoryFuel,
		repoModel.PartCategoryPorthole,
		repoModel.PartCategoryWing,
	}

	for _, category := range categories {
		source := repoModel.Part{Category: category}
		result := converter.PartToModel(source)
		s.Equal(model.PartCategory(category), result.Category)
	}
}
