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

// TestUpdatePart_Success verifies that an existing part can be successfully updated in the repository.
func TestUpdatePart_Success(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()

		ctx          = context.Background()
		originalPart = model.Part{
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

	err := inventoryRepository.AddPart(ctx, originalPart)
	require.NoError(t, err)

	// Create updated version of the part
	updatedPart := originalPart
	updatedPart.Name = gofakeit.Word()
	updatedPart.Description = gofakeit.Sentence(15)
	updatedPart.Price = gofakeit.Float64Range(1001, 2000)
	updatedPart.StockQuantity = int64(gofakeit.IntRange(101, 200))

	// Test
	err = inventoryRepository.UpdatePart(ctx, updatedPart)
	require.NoError(t, err)

	// Verify
	retrievedPart, err := inventoryRepository.GetPart(ctx, originalPart.Uuid)
	require.NoError(t, err)
	require.Equal(t, updatedPart, *retrievedPart)
	require.NotEqual(t, originalPart.Name, retrievedPart.Name)
	require.NotEqual(t, originalPart.Description, retrievedPart.Description)
	require.NotEqual(t, originalPart.Price, retrievedPart.Price)
	require.NotEqual(t, originalPart.StockQuantity, retrievedPart.StockQuantity)
}

// TestUpdatePart_NotFound verifies that attempting to update a non-existent part
func TestUpdatePart_NotFound(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		nonExistentPart     = model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Word(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Float64Range(1, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      model.PartCategory(gofakeit.IntRange(1, 4)),
		}
	)

	// Test
	err := inventoryRepository.UpdatePart(ctx, nonExistentPart)

	// Verify
	require.Error(t, err)
	require.Equal(t, repository.ErrPartWithUUIDNotFound(nonExistentPart.Uuid), err)
}

// TestUpdatePart_EmptyUUID verifies the repository correctly handles update attempts
func TestUpdatePart_EmptyUUID(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		partWithEmptyUUID   = model.Part{
			Uuid:          "",
			Name:          gofakeit.Word(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Float64Range(1, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
		}
	)

	// Test
	err := inventoryRepository.UpdatePart(ctx, partWithEmptyUUID)

	// Verify
	require.Error(t, err)
	require.Equal(t, repository.ErrPartWithUUIDNotFound(""), err)
}

// TestUpdatePart_ConcurrentAccess verifies thread-safe update behavior by performing
// concurrent updates to the same part from multiple goroutines.
func TestUpdatePart_ConcurrentAccess(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		originalPart        = model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Word(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Float64Range(1, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      model.PartCategory(gofakeit.IntRange(1, 4)),
		}
	)

	err := inventoryRepository.AddPart(ctx, originalPart)
	require.NoError(t, err)

	// Test
	var wg sync.WaitGroup
	updateCount := 5
	wg.Add(updateCount)

	for i := 0; i < updateCount; i++ {
		go func(iteration int) {
			defer wg.Done()
			updatedPart := originalPart
			updatedPart.Name = gofakeit.Word()
			updatedPart.StockQuantity = int64(iteration + 1)
			err := inventoryRepository.UpdatePart(ctx, updatedPart)
			require.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Verify
	finalPart, err := inventoryRepository.GetPart(ctx, originalPart.Uuid)
	require.NoError(t, err)
	require.NotEqual(t, originalPart.Name, finalPart.Name)
	require.NotEqual(t, originalPart.StockQuantity, finalPart.StockQuantity)
}

// TestUpdatePart_DefensiveCopy verifies the repository makes defensive copies of parts
func TestUpdatePart_DefensiveCopy(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		originalPart        = model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Word(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Float64Range(1, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
		}
	)

	err := inventoryRepository.AddPart(ctx, originalPart)
	require.NoError(t, err)

	// Create updated part and modify it after update
	updatedPart := originalPart
	updatedPart.Name = "New Name"

	// Test
	err = inventoryRepository.UpdatePart(ctx, updatedPart)
	require.NoError(t, err)

	// Modify the local copy after update
	updatedPart.Name = "Modified After Update"

	// Verify the repository copy wasn't affected
	retrievedPart, err := inventoryRepository.GetPart(ctx, originalPart.Uuid)
	require.NoError(t, err)
	require.Equal(t, "New Name", retrievedPart.Name)
	require.NotEqual(t, updatedPart.Name, retrievedPart.Name)
}
