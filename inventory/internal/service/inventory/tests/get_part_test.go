package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dvln/testify/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/service/inventory"
	"github.com/andredubov/rocket-factory/inventory/internal/service/mocks"
)

// TestGetPart_Success verifies successful retrieval of a part through the service layer.
func TestGetPart_Success(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

		ctx  = context.Background()
		uuid = gofakeit.UUID()

		expectedPart = &model.Part{
			Uuid:          uuid,
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
	inventoryRepository.On("GetPart", mock.Anything, uuid).Return(expectedPart, nil).Once()

	// Test
	retrievedPart, err := inventoryService.GetPart(ctx, uuid)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, retrievedPart)
	assert.Equal(t, expectedPart, retrievedPart)
	assert.Equal(t, uuid, retrievedPart.Uuid)
	assert.NotEqual(t, "", retrievedPart.Name)
	assert.NotEqual(t, "", retrievedPart.Description)
	assert.True(t, retrievedPart.Price > 0)
	assert.True(t, retrievedPart.StockQuantity > 0)
	assert.True(t, len(retrievedPart.Tags) > 0)
	assert.NotNil(t, retrievedPart.Dimensions)
	assert.NotNil(t, retrievedPart.Manufacturer)

	inventoryRepository.AssertExpectations(t)
}

// TestGetPart_NotFoundError verifies proper error handling when requesting a non-existent part.
func TestGetPart_NotFoundError(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

		ctx         = context.Background()
		uuid        = gofakeit.UUID()
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	inventoryRepository.On("GetPart", mock.Anything, uuid).Return(nil, expectedErr).Once()

	// Test
	retrievedPart, err := inventoryService.GetPart(ctx, uuid)

	// Verify
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	require.Nil(t, retrievedPart)

	inventoryRepository.AssertExpectations(t)
}

// TestGetPart_EmptyUUID verifies the service handles empty UUID values appropriately.
func TestGetPart_EmptyUUID(t *testing.T) {
	// Setup
	var (
		inventoryRepository = mocks.NewInventoryRepository(t)
		inventoryService    = inventory.NewService(inventoryRepository)

		ctx         = context.Background()
		emptyUUID   = ""
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	inventoryRepository.On("GetPart", mock.Anything, emptyUUID).Return(nil, expectedErr).Once()

	// Test
	part, err := inventoryService.GetPart(ctx, emptyUUID)

	// Verify
	require.Nil(t, part)
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)

	// Ensure repository was never called
	inventoryRepository.AssertNotCalled(t, "GetPart")
	inventoryRepository.AssertExpectations(t)
}
