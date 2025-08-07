package tests

import (
	"context"
	"fmt"
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

func TestUpdatePart_ConcurrentAccess(t *testing.T) {
	// Настройка
	var (
		inventoryRepository = memory.NewInventoryRepository()
		ctx                 = context.Background()
		originalPart        = model.Part{
			Uuid:          "concurrent-test-uuid",
			Name:          "Исходное название",
			Description:   "Тестовая часть для проверки конкурентного доступа",
			Price:         100.0,
			StockQuantity: 0, // Начинаем с 0
			Category:      model.PartCategoryEngine,
		}
	)

	// Добавляем оригинальную часть в репозиторий
	err := inventoryRepository.AddPart(ctx, originalPart)
	require.NoError(t, err)

	// Количество конкурентных обновлений
	updateCount := 5
	var wg sync.WaitGroup
	wg.Add(updateCount)

	// Канал для сбора результатов обновлений
	results := make(chan int64, updateCount)

	// Запускаем горутины для конкурентных обновлений
	for i := 0; i < updateCount; i++ {
		go func(iteration int) {
			defer wg.Done()
			updatedPart := originalPart
			updatedPart.Name = fmt.Sprintf("Обновлённое-имя-%d", iteration) // Устанавливаем уникальные имена
			updatedPart.StockQuantity = int64(iteration + 1)                // Устанавливаем значения от 1 до 5
			err := inventoryRepository.UpdatePart(ctx, updatedPart)
			require.NoError(t, err)
			results <- updatedPart.StockQuantity
		}(i)
	}

	// Ждём завершения всех горутин
	wg.Wait()
	close(results)

	// Получаем финальную версию части
	finalPart, err := inventoryRepository.GetPart(ctx, originalPart.Uuid)
	require.NoError(t, err)

	// Проверяем, что имя изменилось (должно быть одним из установленных значений)
	require.NotEqual(t, originalPart.Name, finalPart.Name)
	require.Contains(t, []string{
		"Обновлённое-имя-0",
		"Обновлённое-имя-1",
		"Обновлённое-имя-2",
		"Обновлённое-имя-3",
		"Обновлённое-имя-4",
	}, finalPart.Name)

	// Проверяем, что количество на складе соответствует одному из установленных значений
	validValues := make(map[int64]bool)
	for i := 1; i <= updateCount; i++ {
		validValues[int64(i)] = true
	}

	require.True(t, validValues[finalPart.StockQuantity],
		"StockQuantity должен быть одним из значений от 1 до %d, получено %d",
		updateCount, finalPart.StockQuantity)

	// Дополнительная проверка: все отправленные значения должны быть среди validValues
	for quantity := range results {
		require.True(t, validValues[quantity])
	}
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
