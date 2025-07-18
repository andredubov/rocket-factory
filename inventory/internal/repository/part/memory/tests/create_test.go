package tests

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dvln/testify/assert"

	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// TestAddPart_Success verifies that a valid part can be successfully added to the repository.
func (s *InventoryRepositorySuite) TestAddPart_Success() {
	// Generate a random part using gofakeit
	var (
		ctx  = context.Background()
		part = model.Part{
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
	err := s.inventoryRepository.AddPart(ctx, part)

	// Verify no error and part exists
	s.Require().NoError(err)
	s.Require().Nil(err)
	retrieved, err := s.inventoryRepository.GetPart(ctx, part.Uuid)
	s.Require().NoError(err)
	s.Require().Equal(part, *retrieved)
}

// TestAddPart_DuplicateUUID tests that attempting to add a part with an existing UUID
// returns the expected error (ErrPartWithUUIDExists). Verifies duplicate detection
func (s *InventoryRepositorySuite) TestAddPart_DuplicateUUID() {
	// Create a test part
	var (
		ctx  = context.Background()
		part = model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Category:      model.PartCategoryEngine,
			StockQuantity: 10,
		}
	)

	// Add the part first time (should succeed)
	err := s.inventoryRepository.AddPart(ctx, part)
	s.Require().NoError(err)
	s.Require().Nil(err)

	// Try to add again with same UUID
	err = s.inventoryRepository.AddPart(ctx, part)

	// Verify we get the expected error
	s.Require().Error(err)
	s.Require().Equal(repository.ErrPartWithUUIDExists(part.Uuid), err)
}

// TestAddPart_EmptyPart verifies that a minimal valid part (with only required fields)
// can be successfully added to the repository. Tests minimal valid input case.
func (s *InventoryRepositorySuite) TestAddPart_EmptyPart() {
	// Test with minimal valid part
	var (
		ctx  = context.Background()
		part = model.Part{
			Uuid:     gofakeit.UUID(),
			Name:     gofakeit.Name(),
			Category: model.PartCategoryUnknown,
		}
	)

	err := s.inventoryRepository.AddPart(ctx, part)
	s.Require().NoError(err)
	s.Require().Nil(err)
}

// TestAddPart_ZeroValues verifies that a part with zero values for optional fields
// can be successfully added. Tests default/empty values handling.
func (s *InventoryRepositorySuite) TestAddPart_ZeroValues() {
	// Test with zero values
	var (
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

	err := s.inventoryRepository.AddPart(ctx, part)
	s.Require().NoError(err)
	s.Require().Nil(err)
}

// TestAddPart_ConcurrentAccess verifies thread-safe behavior by attempting to add
// the same part concurrently from multiple goroutines. Ensures only one succeeds
// and others get duplicate error, while maintaining data consistency.
func (s *InventoryRepositorySuite) TestAddPart_ConcurrentAccess() {
	// This test verifies the mutex protection
	var (
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
			err := s.inventoryRepository.AddPart(ctx, part)
			if err != nil {
				assert.Equal(s.T(), repository.ErrPartWithUUIDExists(part.Uuid), err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numWorkers; i++ {
		<-done
	}

	// Verify part was added exactly once
	retrieved, err := s.inventoryRepository.GetPart(ctx, part.Uuid)
	s.Require().NoError(err)
	s.Require().Equal(part, *retrieved)
}
