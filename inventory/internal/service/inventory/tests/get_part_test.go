package tests

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// TestGetPart_Success verifies successful retrieval of a part through the service layer.
// Tests that a valid UUID returns the corresponding part with all fields properly converted
// from the repository model to the domain model. Validates correct data transformation.
func (s *InventoryServiceSuite) TestGetPart_Success() {
	// Setup
	var (
		ctx  = context.Background()
		uuid = gofakeit.UUID()

		repoPart = &repoModel.Part{
			Uuid:          uuid,
			Name:          gofakeit.Word(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Float64Range(1, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      repoModel.PartCategory(gofakeit.IntRange(1, 4)),
			Dimensions: repoModel.Dimensions{
				Length: gofakeit.Float64Range(1, 100),
				Width:  gofakeit.Float64Range(1, 100),
				Height: gofakeit.Float64Range(1, 100),
				Weight: gofakeit.Float64Range(1, 100),
			},
			Manufacturer: repoModel.Manufacturer{
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
	s.inventoryRepository.On("GetPart", ctx, uuid).Return(repoPart, nil)

	// Test
	part, err := s.inventoryService.GetPart(ctx, uuid)

	// Verify
	s.Require().NoError(err)
	s.Require().NotNil(part)
	s.Require().Equal(uuid, part.Uuid)
	s.Require().Equal(converter.PartToModel(*repoPart), *part)
}

// TestGetPart_NotFoundError verifies proper error handling when requesting a non-existent part.
// Tests that the repository's ErrPartNotFound error is correctly propagated through
// the service layer while maintaining the original error type.
func (s *InventoryServiceSuite) TestGetPart_NotFoundError() {
	// Setup
	var (
		ctx         = context.Background()
		uuid        = gofakeit.UUID()
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	s.inventoryRepository.On("GetPart", ctx, uuid).Return(nil, expectedErr)

	// Test
	part, err := s.inventoryService.GetPart(ctx, uuid)

	// Verify
	s.Require().Error(err)
	s.Require().Equal(err, expectedErr)
	s.Require().Nil(part)
}

// TestGetPart_EmptyUUID verifies the service handles empty UUID values appropriately.
// Tests that attempting to retrieve a part with an empty UUID returns the expected
// ErrPartNotFound error, validating input sanitization at the service boundary.
func (s *InventoryServiceSuite) TestGetPart_EmptyUUID() {
	// Setup
	var (
		ctx         = context.Background()
		emptyUUID   = ""
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	s.inventoryRepository.On("GetPart", ctx, emptyUUID).Return(nil, expectedErr)

	// Test
	part, err := s.inventoryService.GetPart(ctx, emptyUUID)

	// Verify
	s.Nil(part)
	s.Error(err)
	s.Require().Equal(err, expectedErr)
}
