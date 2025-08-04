package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/service/inventory"
	"github.com/andredubov/rocket-factory/inventory/internal/service/mocks"
)

// TestGetPartList_Success verifies successful retrieval of filtered parts through the service layer.
func TestGetPartList_Success(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

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
	inventoryRepository.On("GetPartList", ctx, filter).Return(parts, nil).Once()

	// Test
	retrivedParts, err := inventoryService.GetPartList(ctx, filter)

	// Verify
	require.NoError(t, err)
	require.Len(t, retrivedParts, len(parts))
	assert.Equal(t, retrivedParts, parts)
	inventoryRepository.AssertExpectations(t)
}

// TestGetPartList_EmptyFilter verifies behavior when querying with an empty filter.
func TestGetPartList_EmptyFilter(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

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
	inventoryRepository.On("GetPartList", ctx, filter).Return(parts, nil)

	// Test
	retrivedParts, err := inventoryService.GetPartList(ctx, filter)

	// Verify
	require.NoError(t, err)
	require.Len(t, retrivedParts, len(parts))
	assert.Equal(t, retrivedParts, parts)
	inventoryRepository.AssertExpectations(t)
}

// TestGetPartList_RepositoryError verifies proper error propagation from the repository.
func TestGetPartList_RepositoryError(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

		ctx    = context.Background()
		filter = model.PartFilter{
			UUIDs: []string{gofakeit.UUID()},
		}
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	inventoryRepository.On("GetPartList", ctx, filter).Return(nil, expectedErr)

	// Test
	parts, err := inventoryService.GetPartList(ctx, filter)

	// Verify
	require.Nil(t, parts)
	require.Error(t, err)
	require.Equal(t, err, expectedErr)
}

// TestGetPartList_EmptyResult verifies correct handling of empty result sets.
func TestGetPartList_EmptyResult(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

		ctx    = context.Background()
		filter = model.PartFilter{
			Tags: []string{"non-existent-tag"},
		}
	)

	// Mock expectations
	inventoryRepository.On("GetPartList", ctx, filter).Return([]model.Part{}, nil)

	// Test
	retrievedParts, err := inventoryService.GetPartList(ctx, filter)

	// Verify
	require.NoError(t, err)
	assert.Empty(t, retrievedParts)
	inventoryRepository.AssertExpectations(t)
}
