package tests

import "github.com/andredubov/rocket-factory/inventory/internal/repository/model"

// TestIsValid_ValidCategories verifies that all defined part categories
// are recognized as valid by the IsValid() method. Tests the complete set
// of known good values including PartCategoryUnknown.
func (s *InventoryRepoModelsSuite) TestIsValid_ValidCategories() {
	validCategories := []model.PartCategory{
		model.PartCategoryUnknown,
		model.PartCategoryEngine,
		model.PartCategoryFuel,
		model.PartCategoryPorthole,
		model.PartCategoryWing,
	}

	for _, category := range validCategories {
		s.True(category.IsValid(), "Category %d should be valid", category)
	}
}

// TestIsValid_InvalidCategories verifies that clearly invalid category values
// (out of defined range) are properly rejected by the IsValid() method.
// Tests negative numbers and values above the maximum defined category.
func (s *InventoryRepoModelsSuite) TestIsValid_InvalidCategories() {
	invalidCategories := []model.PartCategory{-1, 5, 100}

	for _, category := range invalidCategories {
		s.False(category.IsValid(), "Category %d should be invalid", category)
	}
}

// TestIsValid_BoundaryCases specifically tests the boundary values around
// the valid category range. Verifies both edge values (0 and 4) are accepted
// while adjacent invalid values (-1 and 5) are rejected.
func (s *InventoryRepoModelsSuite) TestIsValid_BoundaryCases() {
	boundaryTests := []struct {
		category model.PartCategory
		expected bool
	}{
		{0, true},   // Min boundary (Unknown)
		{4, true},   // Max boundary (Wing)
		{-1, false}, // Below min
		{5, false},  // Above max
	}

	for _, test := range boundaryTests {
		s.Equal(test.expected, test.category.IsValid(),
			"Category %d validation should be %v", test.category, test.expected)
	}
}
