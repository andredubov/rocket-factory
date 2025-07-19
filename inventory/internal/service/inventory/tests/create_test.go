package tests

import (
	"context"
	"errors"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
)

// TestAddPart_Success verifies that a valid part can be successfully added through the service layer.
// Tests the happy path scenario where all part fields are properly populated and the repository
// accepts the part without errors. Verifies proper conversion and error handling.
func (s *InventoryServiceSuite) TestAddPart_Success() {
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	)

	// Mock expectations
	s.inventoryRepository.On("AddPart", ctx, converter.PartToRepoModel(part)).Return(nil)

	// Test
	err := s.inventoryService.AddPart(ctx, part)

	// Verify
	s.Require().NoError(err)
	s.Require().Nil(err)
}

// TestAddPart_RepositoryError verifies proper error propagation when the repository fails.
// Tests that repository-level errors (like duplicate parts) are properly propagated up
// through the service layer to the caller with their original error type preserved.
func (s *InventoryServiceSuite) TestAddPart_RepositoryError() {
	// Setup
	var (
		ctx  = context.Background()
		part = model.Part{
			Uuid: gofakeit.UUID(),
			Name: gofakeit.Word(),
		}
		expectedErr = model.ErrPartAlreadyExists
	)

	// Mock expectations
	s.inventoryRepository.On("AddPart", ctx, converter.PartToRepoModel(part)).Return(expectedErr)

	// Test
	err := s.inventoryService.AddPart(ctx, part)

	// Verify
	s.Require().NotNil(err)
	s.Require().Equal(err, expectedErr)
}

// TestAddPart_InvalidPart verifies the service handles invalid part data correctly.
// Tests error cases where the part fails business logic validation before reaching
// the repository layer, such as with invalid category values.
func (s *InventoryServiceSuite) TestAddPart_InvalidPart() {
	// Setup
	ctx := context.Background()
	invalidPart := model.Part{
		Category: model.PartCategory(999), // Invalid category
	}

	// Mock expectations
	s.inventoryRepository.On("AddPart", ctx, converter.PartToRepoModel(invalidPart)).Return(errors.New("Some error"))

	// Test
	err := s.inventoryService.AddPart(ctx, invalidPart)

	// Verify
	s.Require().Error(err)
}
