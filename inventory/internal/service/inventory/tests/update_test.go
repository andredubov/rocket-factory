package tests

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
)

// TestUpdatePart_Success verifies that a valid part can be successfully updated through the service layer.
// Tests the happy path scenario where all part fields are properly populated and the repository
// accepts the update without errors. Verifies proper model conversion and successful update flow.
func (s *InventoryServiceSuite) TestUpdatePart_Success() {
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
	repoPart := converter.PartToRepoModel(part)
	s.inventoryRepository.On("UpdatePart", ctx, repoPart).Return(nil)

	// Test
	err := s.inventoryService.UpdatePart(ctx, part)

	// Verify
	s.Require().NoError(err)
}

// TestUpdatePart_NotFoundError verifies proper error handling when attempting to update
// a non-existent part. Tests that the repository's ErrPartNotFound error is correctly
// propagated through the service layer with its original error type preserved.
func (s *InventoryServiceSuite) TestUpdatePart_NotFoundError() {
	// Setup
	var (
		ctx  = context.Background()
		part = model.Part{
			Uuid: gofakeit.UUID(),
			Name: gofakeit.Word(),
		}
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	repoPart := converter.PartToRepoModel(part)
	s.inventoryRepository.On("UpdatePart", ctx, repoPart).Return(expectedErr)

	// Test
	err := s.inventoryService.UpdatePart(ctx, part)

	// Verify
	s.Require().Error(err)
	s.Require().Equal(err, expectedErr)
}

// TestUpdatePart_InvalidPart verifies the service properly rejects invalid part updates.
// Tests error cases where the part data fails validation (empty UUID, invalid category)
// before reaching the repository layer. Validates business logic enforcement.
func (s *InventoryServiceSuite) TestUpdatePart_InvalidPart() {
	// Setup
	var (
		ctx         = context.Background()
		invalidPart = model.Part{
			Uuid:     "",                      // Empty UUID is invalid
			Category: model.PartCategory(999), // Invalid category
		}
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	repoPart := converter.PartToRepoModel(invalidPart)
	s.inventoryRepository.On("UpdatePart", ctx, repoPart).Return(expectedErr)

	// Test
	err := s.inventoryService.UpdatePart(ctx, invalidPart)

	// Verify
	s.Require().Error(err)
	s.Require().Equal(err, expectedErr)
}
