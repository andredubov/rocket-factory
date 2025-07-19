package tests

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
)

// TestDeletePart_Success verifies that a part can be successfully deleted when it exists.
// Tests the happy path scenario where the repository successfully deletes the part
// and no error is returned to the caller.
func (s *InventoryServiceSuite) TestDeletePart_Success() {
	// Setup
	var (
		ctx  = context.Background()
		uuid = gofakeit.UUID()
	)

	// Mock expectations
	s.inventoryRepository.On("DeletePart", ctx, uuid).Return(nil)

	// Test
	err := s.inventoryService.DeletePart(ctx, uuid)

	// Verify
	s.Require().NoError(err)
}

// TestDeletePart_NotFoundError verifies the service properly handles cases where
// the part to be deleted doesn't exist. Tests that the repository's ErrPartNotFound
// error is correctly propagated through the service layer.
func (s *InventoryServiceSuite) TestDeletePart_NotFoundError() {
	// Setup
	var (
		ctx         = context.Background()
		uuid        = gofakeit.UUID()
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	s.inventoryRepository.On("DeletePart", ctx, uuid).Return(expectedErr)

	// Test
	err := s.inventoryService.DeletePart(ctx, uuid)

	// Verify
	s.Require().NotNil(err)
	s.Require().Equal(err, expectedErr)
}

// TestDeletePart_EmptyUUID verifies the service handles empty UUID values correctly.
// Tests that attempting to delete a part with an empty UUID returns the expected
// ErrPartNotFound error, validating input sanitization.
func (s *InventoryServiceSuite) TestDeletePart_EmptyUUID() {
	// Setup
	var (
		ctx         = context.Background()
		emptyUUID   = ""
		expectedErr = model.ErrPartNotFound
	)

	// Mock expectations
	s.inventoryRepository.On("DeletePart", ctx, emptyUUID).Return(expectedErr)

	// Test
	err := s.inventoryService.DeletePart(ctx, emptyUUID)

	// Verify
	s.Require().Error(err)
	s.Require().Equal(err, expectedErr)
}
