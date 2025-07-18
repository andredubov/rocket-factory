package tests

import (
	"context"
	"sync"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// TestGetPart_Success verifies that a part can be successfully retrieved from the repository
// when it exists. It first adds a test part with all fields populated, then verifies
// the retrieved part matches exactly what was stored.
func (s *InventoryRepositorySuite) TestGetPart_Success() {
	// Setup
	var (
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

	err := s.inventoryRepository.AddPart(ctx, expectedPart)
	s.Require().NoError(err)
	s.Require().Nil(err)

	// Test
	actualPart, err := s.inventoryRepository.GetPart(ctx, expectedPart.Uuid)

	// Verify
	s.Require().NoError(err)
	s.Require().NotNil(actualPart)
	s.Require().Equal(expectedPart, *actualPart)
}

// TestGetPart_NotFound verifies the repository correctly returns ErrPartWithUUIDNotFound
// when attempting to get a part with a non-existent UUID. Tests error handling for missing parts.
func (s *InventoryRepositorySuite) TestGetPart_NotFound() {
	// Setup
	var (
		ctx             = context.Background()
		nonExistentUUID = gofakeit.UUID()
	)

	// Test
	part, err := s.inventoryRepository.GetPart(ctx, nonExistentUUID)

	// Verify
	s.Require().Error(err)
	s.Require().Equal(repository.ErrPartWithUUIDNotFound(nonExistentUUID), err)
	s.Require().Nil(part)
}

// TestGetPart_EmptyUUID verifies the repository handles empty UUID strings properly
// by returning ErrPartWithUUIDNotFound. Tests edge case for invalid input.
func (s *InventoryRepositorySuite) TestGetPart_EmptyUUID() {
	// Setup
	var (
		ctx       = context.Background()
		emptyUUID = ""
	)

	// Test
	part, err := s.inventoryRepository.GetPart(ctx, emptyUUID)

	// Verify
	s.Require().Error(err)
	s.Require().Equal(repository.ErrPartWithUUIDNotFound(emptyUUID), err)
	s.Require().Nil(part)
}

// TestGetPart_ConcurrentAccess verifies thread-safe read behavior by concurrently
// accessing the same part from multiple goroutines. Ensures all goroutines receive
// the correct part data without errors, validating read consistency under concurrency.
func (s *InventoryRepositorySuite) TestGetPart_ConcurrentAccess() {
	// Setup
	var (
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

	err := s.inventoryRepository.AddPart(ctx, expectedPart)
	s.Require().NoError(err)
	s.Require().Nil(err)

	// Test
	var wg sync.WaitGroup
	results := make(chan *model.Part, 10)
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			part, err := s.inventoryRepository.GetPart(ctx, expectedPart.Uuid)
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
		s.Require().Equal(expectedPart, *part)
	}

	for err := range errors {
		s.Require().NoError(err) // Should never receive errors in this case
		s.Require().Nil(err)
	}
}
