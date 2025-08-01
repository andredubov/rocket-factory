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

// TestDeletePart_Success verifies that a part can be successfully deleted when it exists.
func TestDeletePart_Success(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

		ctx  = context.Background()
		uuid = gofakeit.UUID()
	)

	// Mock expectations
	inventoryRepository.On("DeletePart", ctx, uuid).Return(nil).Once()

	// Test
	err := inventoryService.DeletePart(ctx, uuid)

	// Verify
	require.NoError(t, err)
	inventoryRepository.AssertExpectations(t)
	inventoryRepository.AssertCalled(t, "DeletePart", ctx, uuid)
	inventoryRepository.AssertNumberOfCalls(t, "DeletePart", 1)
}

// TestDeletePart_NotFoundError verifies the service properly handles cases where
func TestDeletePart_NotFoundError(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

		ctx         = context.Background()
		uuid        = gofakeit.UUID()
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	inventoryRepository.On("DeletePart", ctx, uuid).Return(expectedErr).Once()

	// Test
	err := inventoryService.DeletePart(ctx, uuid)

	// Verify
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	inventoryRepository.AssertExpectations(t)
}

// TestDeletePart_EmptyUUID verifies the service handles empty UUID values correctly.
func TestDeletePart_EmptyUUID(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

		ctx         = context.Background()
		emptyUUID   = ""
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	inventoryRepository.On("DeletePart", ctx, emptyUUID).Return(expectedErr).Once()

	// Test
	err := inventoryService.DeletePart(ctx, emptyUUID)

	// Verify
	require.Error(t, err)
	require.Equal(t, err, expectedErr)
	inventoryRepository.AssertExpectations(t)
}
