package tests

import (
	"context"
	"sync"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// TestUpdatePart_Success verifies that an existing part can be successfully updated in the repository.
func (s *InventoryRepositorySuite) TestUpdatePart_Success() {
	// Setup
	var (
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

	err := s.inventoryRepository.AddPart(ctx, originalPart)
	s.Require().NoError(err)
	s.Require().Nil(err)

	// Create updated version of the part
	updatedPart := originalPart
	updatedPart.Name = gofakeit.Word()
	updatedPart.Description = gofakeit.Sentence(15)
	updatedPart.Price = gofakeit.Float64Range(1001, 2000)
	updatedPart.StockQuantity = int64(gofakeit.IntRange(101, 200))

	// Test
	err = s.inventoryRepository.UpdatePart(ctx, updatedPart)
	s.Require().NoError(err)
	s.Require().Nil(err)

	// Verify
	retrievedPart, err := s.inventoryRepository.GetPart(ctx, originalPart.Uuid)
	s.Require().NoError(err)
	s.Require().Nil(err)
	s.Require().Equal(updatedPart, *retrievedPart)
	s.Require().NotEqual(originalPart.Name, retrievedPart.Name)
	s.Require().NotEqual(originalPart.Description, retrievedPart.Description)
	s.Require().NotEqual(originalPart.Price, retrievedPart.Price)
	s.Require().NotEqual(originalPart.StockQuantity, retrievedPart.StockQuantity)
}

// TestUpdatePart_NotFound verifies that attempting to update a non-existent part
// returns the expected ErrPartWithUUIDNotFound error.
func (s *InventoryRepositorySuite) TestUpdatePart_NotFound() {
	// Setup
	ctx := context.Background()
	nonExistentPart := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Word(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Float64Range(1, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      model.PartCategory(gofakeit.IntRange(1, 4)),
	}

	// Test
	err := s.inventoryRepository.UpdatePart(ctx, nonExistentPart)

	// Verify
	s.Require().Error(err)
	s.Require().Equal(repository.ErrPartWithUUIDNotFound(nonExistentPart.Uuid), err)
}

// TestUpdatePart_EmptyUUID verifies the repository correctly handles update attempts
// with empty UUID strings by returning ErrPartWithUUIDNotFound.
func (s *InventoryRepositorySuite) TestUpdatePart_EmptyUUID() {
	// Setup
	ctx := context.Background()
	partWithEmptyUUID := model.Part{
		Uuid:          "",
		Name:          gofakeit.Word(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Float64Range(1, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
	}

	// Test
	err := s.inventoryRepository.UpdatePart(ctx, partWithEmptyUUID)

	// Verify
	s.Require().Error(err)
	s.Require().Equal(repository.ErrPartWithUUIDNotFound(""), err)
}

// TestUpdatePart_ConcurrentAccess verifies thread-safe update behavior by performing
// concurrent updates to the same part from multiple goroutines.
func (s *InventoryRepositorySuite) TestUpdatePart_ConcurrentAccess() {
	// Setup
	ctx := context.Background()
	originalPart := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Word(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Float64Range(1, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      model.PartCategory(gofakeit.IntRange(1, 4)),
	}

	err := s.inventoryRepository.AddPart(ctx, originalPart)
	s.Require().NoError(err)
	s.Require().Nil(err)

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
			err := s.inventoryRepository.UpdatePart(ctx, updatedPart)
			s.Require().NoError(err)
		}(i)
	}

	wg.Wait()

	// Verify
	finalPart, err := s.inventoryRepository.GetPart(ctx, originalPart.Uuid)
	s.Require().NoError(err)
	s.Require().Nil(err)
	s.Require().NotEqual(originalPart.Name, finalPart.Name)
	s.Require().NotEqual(originalPart.StockQuantity, finalPart.StockQuantity)
}

// TestUpdatePart_DefensiveCopy verifies the repository makes defensive copies of parts
// during updates. Tests that modifications to the local part variable after update
// don't affect the stored version, ensuring data integrity.
func (s *InventoryRepositorySuite) TestUpdatePart_DefensiveCopy() {
	// Setup
	ctx := context.Background()
	originalPart := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Word(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Float64Range(1, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
	}

	err := s.inventoryRepository.AddPart(ctx, originalPart)
	s.Require().NoError(err)
	s.Require().Nil(err)

	// Create updated part and modify it after update
	updatedPart := originalPart
	updatedPart.Name = "New Name"

	// Test
	err = s.inventoryRepository.UpdatePart(ctx, updatedPart)
	s.Require().NoError(err)
	s.Require().Nil(err)

	// Modify the local copy after update
	updatedPart.Name = "Modified After Update"

	// Verify the repository copy wasn't affected
	retrievedPart, err := s.inventoryRepository.GetPart(ctx, originalPart.Uuid)
	s.Require().NoError(err)
	s.Require().Nil(err)
	s.Require().Equal("New Name", retrievedPart.Name)
	s.Require().NotEqual(updatedPart.Name, retrievedPart.Name)
}
