package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/memory"
)

// TestGetPart_Success verifies that a part can be successfully retrieved from the repository
// when it exists. It first adds a test part with all fields populated, then verifies
// the retrieved part matches exactly what was stored.
func TestGetPart_Success(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()

		ctx          = context.Background()
		expectedPart = model.Part{
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

	err := inventoryRepository.AddPart(ctx, expectedPart)
	require.NoError(t, err)

	// Test
	actualPart, err := inventoryRepository.GetPart(ctx, expectedPart.Uuid)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, actualPart)
	require.Equal(t, expectedPart, *actualPart)
}

// TestGetPart_NotFound verifies the repository correctly returns ErrPartWithUUIDNotFound
// when attempting to get a part with a non-existent UUID. Tests error handling for missing parts.
func TestGetPart_NotFound(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		nonExistentUUID     = gofakeit.UUID()
	)

	// Test
	part, err := inventoryRepository.GetPart(ctx, nonExistentUUID)

	// Verify
	require.Error(t, err)
	require.Equal(t, repository.ErrPartWithUUIDNotFound(nonExistentUUID), err)
	require.Nil(t, part)
}

// TestGetPart_EmptyUUID verifies the repository handles empty UUID strings properly
// by returning ErrPartWithUUIDNotFound. Tests edge case for invalid input.
func TestGetPart_EmptyUUID(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		emptyUUID           = ""
	)

	// Test
	part, err := inventoryRepository.GetPart(ctx, emptyUUID)

	// Verify
	require.Error(t, err)
	require.Equal(t, repository.ErrPartWithUUIDNotFound(emptyUUID), err)
	require.Nil(t, part)
}

// TestGetPart_ConcurrentAccess verifies thread-safe read behavior by concurrently
// accessing the same part from multiple goroutines. Ensures all goroutines receive
// the correct part data without errors, validating read consistency under concurrency.
func TestGetPart_ConcurrentAccess(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()

		ctx          = context.Background()
		expectedPart = model.Part{
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

	err := inventoryRepository.AddPart(ctx, expectedPart)
	require.NoError(t, err)

	// Test
	var wg sync.WaitGroup
	results := make(chan *model.Part, 10)
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			part, err := inventoryRepository.GetPart(ctx, expectedPart.Uuid)
			if err != nil {
				errors <- err
				return
			}
			results <- part
		}()
	}

	wg.Wait()
	close(results)
	close(errors)

	// Verify
	for part := range results {
		require.Equal(t, expectedPart, *part)
	}

	for err := range errors {
		require.NoError(t, err) // Should never receive errors in this case
	}
}
