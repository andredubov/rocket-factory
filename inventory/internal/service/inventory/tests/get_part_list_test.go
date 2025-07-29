package tests

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
)

// TestGetPartList_Success verifies successful retrieval of filtered parts through the service layer.
// Tests that parts matching the filter criteria (UUIDs and categories) are properly returned
// after conversion from repository to domain models. Validates correct filtering and data transformation.
func (s *InventoryServiceSuite) TestGetPartList_Success() {
	// Setup
	var (
		ctx    = context.Background()
		filter = model.PartFilter{
			UUIDs: []string{gofakeit.UUID()},
			Categories: []model.PartCategory{
				model.PartCategoryEngine,
				model.PartCategoryFuel,
			},
		}

		parts = []model.Part{
			{
				Uuid:     gofakeit.UUID(),
				Name:     "Engine Part",
				Category: model.PartCategoryEngine,
				Manufacturer: model.Manufacturer{
					Country: "USA",
				},
			},
			{
				Uuid:     gofakeit.UUID(),
				Name:     "Fuel Part",
				Category: model.PartCategoryFuel,
				Manufacturer: model.Manufacturer{
					Country: "Germany",
				},
			},
		}
	)

	// Mock expectations
	s.inventoryRepository.On("GetPartList", ctx, filter).Return(parts, nil)

	// Test
	retrivedParts, err := s.inventoryService.GetPartList(ctx, filter)

	// Verify
	s.Require().NoError(err)
	s.Require().Len(retrivedParts, 2)
	s.Require().Equal(retrivedParts[0], parts[0])
	s.Require().Equal(retrivedParts[1], parts[1])
}

// TestGetPartList_EmptyFilter verifies behavior when querying with an empty filter.
// Tests that the service returns all available parts when no filter criteria is specified,
// validating default behavior for unfiltered requests.
func (s *InventoryServiceSuite) TestGetPartList_EmptyFilter() {
	// Setup
	var (
		ctx    = context.Background()
		filter = model.PartFilter{} // Empty filter
		parts  = []model.Part{
			{
				Uuid: gofakeit.UUID(),
				Name: "Test Part",
			},
		}
	)

	// Mock expectations
	s.inventoryRepository.On("GetPartList", ctx, filter).Return(parts, nil)

	// Test
	retrivedParts, err := s.inventoryService.GetPartList(ctx, filter)

	// Verify
	s.Require().NoError(err)
	s.Require().Len(parts, 1)
	s.Require().Equal(retrivedParts[0], parts[0])
}

// TestGetPartList_RepositoryError verifies proper error propagation from the repository.
// Tests that repository-level errors (like ErrPartNotFound) are correctly propagated
// through the service layer while maintaining the original error type.
func (s *InventoryServiceSuite) TestGetPartList_RepositoryError() {
	// Setup
	var (
		ctx    = context.Background()
		filter = model.PartFilter{
			UUIDs: []string{gofakeit.UUID()},
		}
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	s.inventoryRepository.On("GetPartList", ctx, filter).Return(nil, expectedErr)

	// Test
	parts, err := s.inventoryService.GetPartList(ctx, filter)

	// Verify
	s.Require().Nil(parts)
	s.Require().Error(err)
	s.Require().Equal(err, expectedErr)
}

// TestGetPartList_EmptyResult verifies correct handling of empty result sets.
// Tests that the service properly returns an empty slice (not nil) when no parts
// match the filter criteria, without returning an error.
func (s *InventoryServiceSuite) TestGetPartList_EmptyResult() {
	// Setup
	var (
		ctx    = context.Background()
		filter = model.PartFilter{
			Tags: []string{"non-existent-tag"},
		}
	)

	// Mock expectations
	s.inventoryRepository.On("GetPartList", ctx, filter).Return([]model.Part{}, nil)

	// Test
	parts, err := s.inventoryService.GetPartList(ctx, filter)

	// Verify
	s.Require().NoError(err)
	s.Require().Empty(parts)
}
