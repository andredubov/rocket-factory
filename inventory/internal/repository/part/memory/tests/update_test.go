package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/memory"
)

func TestUpdatePart_Success(t *testing.T) {
	// Настройка
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()

		// Создаём часть с известными значениями
		originalPart = model.Part{
			Uuid:          "test-uuid-123",
			Name:          "Исходное название",
			Description:   "Исходное описание",
			Price:         100.50,
			StockQuantity: 5,
			Category:      model.PartCategoryEngine,
			Dimensions: model.Dimensions{
				Length: 10.0,
				Width:  5.0,
				Height: 2.0,
				Weight: 1.5,
			},
			Manufacturer: model.Manufacturer{
				Name:    "Оригинальный производитель",
				Country: "Россия",
				Website: "http://example.com",
			},
			Tags:      []string{"тег1", "тег2"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	)

	// Добавляем оригинальную часть в репозиторий
	err := inventoryRepository.AddPart(ctx, originalPart)
	require.NoError(t, err)

	// Создаём обновлённую версию части
	updatedPart := originalPart
	updatedPart.Name = "Обновлённое название"
	updatedPart.Description = "Новое подробное описание"
	updatedPart.Price = 200.75
	updatedPart.StockQuantity = 10
	updatedPart.Tags = []string{"новый_тег"}

	// Выполняем обновление
	err = inventoryRepository.UpdatePart(ctx, updatedPart)
	require.NoError(t, err)

	// Проверяем результаты
	retrievedPart, err := inventoryRepository.GetPart(ctx, originalPart.Uuid)
	require.NoError(t, err)

	// Проверяем обновлённые поля
	require.Equal(t, "Обновлённое название", retrievedPart.Name)
	require.Equal(t, "Новое подробное описание", retrievedPart.Description)
	require.Equal(t, 200.75, retrievedPart.Price)
	require.Equal(t, int64(10), retrievedPart.StockQuantity)
	require.Equal(t, []string{"новый_тег"}, retrievedPart.Tags)

	// Проверяем неизменившиеся поля
	require.Equal(t, originalPart.Uuid, retrievedPart.Uuid)
	require.Equal(t, originalPart.Category, retrievedPart.Category)
	require.Equal(t, originalPart.Dimensions, retrievedPart.Dimensions)
	require.Equal(t, originalPart.Manufacturer, retrievedPart.Manufacturer)
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
