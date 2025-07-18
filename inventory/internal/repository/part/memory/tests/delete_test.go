package tests

import (
	"context"
	"sync"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// TestDeletePart_Success verifies that an existing part can be successfully deleted from the repository.
// It first adds a test part, then deletes it, and finally verifies the part no longer exists.
func (s *InventoryRepositorySuite) TestDeletePart_Success() {
	// Setup
	var (
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

	err := s.inventoryRepository.AddPart(ctx, part)
	s.Require().NoError(err)
	s.Require().Nil(err)

	// Test
	err = s.inventoryRepository.DeletePart(ctx, part.Uuid)

	// Verify
	s.Require().NoError(err)
	s.Require().Nil(err)

	retrived, err := s.inventoryRepository.GetPart(ctx, part.Uuid)
	s.Require().Error(err)
	s.Require().Equal(err, repository.ErrPartWithUUIDNotFound(part.Uuid))
	s.Require().Nil(retrived)
}

// TestDeletePart_NotFound verifies that attempting to delete a non-existent part
// returns the expected ErrPartWithUUIDNotFound error. Tests error handling for missing parts.
func (s *InventoryRepositorySuite) TestDeletePart_NotFound() {
	// Setup
	var (
		ctx             = context.Background()
		nonExistentUUID = gofakeit.UUID()
	)

	// Test
	err := s.inventoryRepository.DeletePart(ctx, nonExistentUUID)

	// Verify
	s.Require().Error(err)
	s.Require().Equal(err, repository.ErrPartWithUUIDNotFound(nonExistentUUID))
}

// TestDeletePart_Concurrent verifies thread-safe deletion behavior by attempting to delete
// the same part concurrently from multiple goroutines. Ensures only one deletion succeeds
// while others receive appropriate not-found errors, maintaining data consistency.
func (s *InventoryRepositorySuite) TestDeletePart_Concurrent() {
	// Setup
	var (
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

	err := s.inventoryRepository.AddPart(ctx, part)
	s.Require().NoError(err)
	s.Require().Nil(err)

	// Test
	var wg sync.WaitGroup
	wg.Add(2)

	var err1, err2 error
	go func() {
		defer wg.Done()
		err1 = s.inventoryRepository.DeletePart(context.Background(), part.Uuid)
	}()

	go func() {
		defer wg.Done()
		err2 = s.inventoryRepository.DeletePart(context.Background(), part.Uuid)
	}()

	wg.Wait()

	// Verify
	s.Assert().True((err1 == nil && err2 != nil) || (err1 != nil && err2 == nil),
		"Exactly one deletion should succeed")

	if err1 == nil {
		s.Require().Error(err2)
		s.Require().Equal(err2, repository.ErrPartWithUUIDNotFound(part.Uuid))
	} else {
		s.Require().Error(err1)
		s.Require().Equal(err1, repository.ErrPartWithUUIDNotFound(part.Uuid))
	}

	retrived, err := s.inventoryRepository.GetPart(ctx, part.Uuid)
	s.Require().Error(err)
	s.Require().Equal(err, repository.ErrPartWithUUIDNotFound(part.Uuid))
	s.Require().Nil(retrived)
}
