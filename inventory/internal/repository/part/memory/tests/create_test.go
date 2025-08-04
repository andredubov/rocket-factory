package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/memory"
)

// TestAddPart_Success verifies that a valid part can be successfully added to the repository.
func TestAddPart_Success(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		part                = model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
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
			Tags:     []string{gofakeit.Word(), gofakeit.Word()},
			Metadata: make(map[string]model.Value),
		}
	)

	// Add the part
	err := inventoryRepository.AddPart(ctx, part)

	// Verify no error and part exists
	require.NoError(t, err)
	retrieved, err := inventoryRepository.GetPart(ctx, part.Uuid)
	require.NoError(t, err)
	require.Equal(t, part, *retrieved)
}

// TestAddPart_DuplicateUUID tests that attempting to add a part with an existing UUID
func TestAddPart_DuplicateUUID(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()

		ctx  = context.Background()
		part = model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Category:      model.PartCategoryEngine,
			StockQuantity: 10,
		}
	)

	// Add the part first time (should succeed)
	err := inventoryRepository.AddPart(ctx, part)
	require.NoError(t, err)
	require.Nil(t, err)

	// Try to add again with same UUID
	err = inventoryRepository.AddPart(ctx, part)

	// Verify we get the expected error
	require.Error(t, err)
	require.Equal(t, repository.ErrPartWithUUIDExists(part.Uuid), err)
}

// TestAddPart_EmptyPart verifies that a minimal valid part (with only required fields)
// can be successfully added to the repository. Tests minimal valid input case.
func TestAddPart_EmptyPart(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()

		ctx  = context.Background()
		part = model.Part{
			Uuid:     gofakeit.UUID(),
			Name:     gofakeit.Name(),
			Category: model.PartCategoryUnknown,
		}
	)

	// Test
	err := inventoryRepository.AddPart(ctx, part)

	// Verify
	require.NoError(t, err)
	require.Nil(t, err)
}

// TestAddPart_ZeroValues verifies that a part with zero values for optional fields
// can be successfully added. Tests default/empty values handling.
func TestAddPart_ZeroValues(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()

		ctx  = context.Background()
		part = model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   "",
			Price:         0,
			StockQuantity: 0,
			Category:      model.PartCategoryUnknown,
			Dimensions:    model.Dimensions{},
			Manufacturer:  model.Manufacturer{},
			Tags:          nil,
			Metadata:      nil,
		}
	)

	// Test
	err := inventoryRepository.AddPart(ctx, part)

	// Verify
	require.NoError(t, err)
}

// TestAddPart_ConcurrentAccess verifies thread-safe behavior by attempting to add
func TestAddPart_ConcurrentAccess(t *testing.T) {
	// Setup
	var (
		inventoryRepository = memory.NewInventoryRepository()

		ctx  = context.Background()
		part = model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Category:      model.PartCategoryWing,
			StockQuantity: 5,
		}
		// Use multiple goroutines to test concurrent access
		done       = make(chan bool)
		numWorkers = 10
	)

	for i := 0; i < numWorkers; i++ {
		go func() {
			err := inventoryRepository.AddPart(ctx, part)
			if err != nil {
				assert.Equal(t, repository.ErrPartWithUUIDExists(part.Uuid), err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numWorkers; i++ {
		<-done
	}

	// Verify part was added exactly once
	retrievedPart, err := inventoryRepository.GetPart(ctx, part.Uuid)
	require.NoError(t, err)
	require.Equal(t, part, *retrievedPart)
}
