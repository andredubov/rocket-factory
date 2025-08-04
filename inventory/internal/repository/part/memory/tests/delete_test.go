package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/memory"
)

// TestDeletePart_Success verifies that an existing part can be successfully deleted from the repository.
func TestDeletePart_Success(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()

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
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}
	)

	err := inventoryRepository.AddPart(ctx, part)
	require.NoError(t, err)
	require.Nil(t, err)

	// Test
	err = inventoryRepository.DeletePart(ctx, part.Uuid)

	// Verify
	require.NoError(t, err)
	retrived, err := inventoryRepository.GetPart(ctx, part.Uuid)
	require.Error(t, err)
	require.Equal(t, err, repository.ErrPartWithUUIDNotFound(part.Uuid))
	require.Nil(t, retrived)
}

// TestDeletePart_NotFound verifies that attempting to delete a non-existent part
func TestDeletePart_NotFound(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		nonExistentUUID     = gofakeit.UUID()
	)

	// Test
	err := inventoryRepository.DeletePart(ctx, nonExistentUUID)

	// Verify
	require.Error(t, err)
	require.Equal(t, err, repository.ErrPartWithUUIDNotFound(nonExistentUUID))
}

// TestDeletePart_Concurrent verifies thread-safe deletion behavior by attempting to delete
func TestDeletePart_Concurrent(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()

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
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}
	)

	err := inventoryRepository.AddPart(ctx, part)
	require.NoError(t, err)

	// Test
	var wg sync.WaitGroup
	wg.Add(2)

	var err1, err2 error
	go func() {
		defer wg.Done()
		err1 = inventoryRepository.DeletePart(context.Background(), part.Uuid)
	}()

	go func() {
		defer wg.Done()
		err2 = inventoryRepository.DeletePart(context.Background(), part.Uuid)
	}()

	wg.Wait()

	// Verify
	assert.True(t, (err1 == nil && err2 != nil) || (err1 != nil && err2 == nil), "Exactly one deletion should succeed")

	if err1 == nil {
		require.Error(t, err2)
		require.Equal(t, err2, repository.ErrPartWithUUIDNotFound(part.Uuid))
	} else {
		require.Error(t, err1)
		require.Equal(t, err1, repository.ErrPartWithUUIDNotFound(part.Uuid))
	}

	retrived, err := inventoryRepository.GetPart(ctx, part.Uuid)
	require.Error(t, err)
	require.Equal(t, err, repository.ErrPartWithUUIDNotFound(part.Uuid))
	require.Nil(t, retrived)
}
