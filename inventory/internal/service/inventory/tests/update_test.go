package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/service/inventory"
	"github.com/andredubov/rocket-factory/inventory/internal/service/mocks"
)

// TestUpdatePart_Success verifies that a valid part can be successfully updated through the service layer.
func TestUpdatePart_Success(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

		ctx  = context.Background()
		part = model.Part{
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
			Tags:      []string{gofakeit.Word(), gofakeit.Word()},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	)

	// Mock expectations
	inventoryRepository.On("UpdatePart", ctx, part).Return(nil).Once()

	// Test
	err := inventoryService.UpdatePart(ctx, part)

	// Verify
	require.NoError(t, err)
	inventoryRepository.AssertExpectations(t)
}

// TestUpdatePart_NotFoundError verifies proper error handling when attempting to update
func TestUpdatePart_NotFoundError(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

		ctx  = context.Background()
		part = model.Part{
			Uuid: gofakeit.UUID(),
			Name: gofakeit.Word(),
		}
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	inventoryRepository.On("UpdatePart", ctx, part).Return(expectedErr)

	// Test
	err := inventoryService.UpdatePart(ctx, part)

	// Verify
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	inventoryRepository.AssertExpectations(t)
}

// TestUpdatePart_InvalidPart verifies the service properly rejects invalid part updates.
func TestUpdatePart_InvalidPart(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

		ctx         = context.Background()
		invalidPart = model.Part{
			Uuid:     "",                      // Empty UUID is invalid
			Category: model.PartCategory(999), // Invalid category
		}
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	inventoryRepository.On("UpdatePart", ctx, invalidPart).Return(expectedErr).Once()

	// Test
	err := inventoryService.UpdatePart(ctx, invalidPart)

	// Verify
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	inventoryRepository.AssertNotCalled(t, "UpdatePart")
	inventoryRepository.AssertExpectations(t)
}
