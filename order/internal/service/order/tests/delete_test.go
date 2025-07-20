package tests

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// TestDeleteOrder_Success verifies successful order deletion through the service layer.
// Tests that valid order UUIDs are properly processed and deleted from the repository
// without errors, confirming the happy path scenario.
func (s *OrdersServiceSuite) TestDeleteOrder_Success() {
	// Arrange
	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
	)

	s.ordersRepository.On("DeleteOrder", ctx, orderUUID).Return(nil)

	// Act
	err := s.ordersService.DeleteOrder(ctx, orderUUID)

	// Assert
	s.NoError(err)
}

// TestDeleteOrder_NotFound verifies proper handling of non-existent order deletions.
// Tests that attempts to delete missing orders return the expected ErrOrderNotFound
// error, maintaining proper error handling semantics.
func (s *OrdersServiceSuite) TestDeleteOrder_NotFound() {
	// Arrange
	var (
		ctx         = context.Background()
		orderUUID   = uuid.New()
		expectedErr = model.ErrOrderNotFound
	)

	s.ordersRepository.On("DeleteOrder", ctx, orderUUID).Return(expectedErr)

	// Act
	err := s.ordersService.DeleteOrder(ctx, orderUUID)

	// Assert
	s.Error(err)
	s.Equal(err, expectedErr)
}

// TestDeleteOrder_WithGofakeit performs randomized testing of order deletion.
// Generates multiple test cases with random success/failure conditions to verify
// the service handles both successful deletions and error cases correctly.
func (s *OrdersServiceSuite) TestDeleteOrder_WithGofakeit() {
	// Generate 5 random test cases
	for i := 0; i < 5; i++ {
		s.Run(gofakeit.BeerName(), func() {
			// Arrange
			ctx := context.Background()
			orderUUID := uuid.New()

			// Randomly decide if test should succeed or fail
			if gofakeit.Bool() {
				// Success case
				s.ordersRepository.On("DeleteOrder", ctx, orderUUID).Return(nil)

				// Act
				err := s.ordersService.DeleteOrder(ctx, orderUUID)

				// Assert
				s.NoError(err)
			} else {
				// Failure case
				expectedErr := model.ErrOrderNotFound
				s.ordersRepository.On("DeleteOrder", ctx, orderUUID).Return(expectedErr)

				// Act
				err := s.ordersService.DeleteOrder(ctx, orderUUID)

				// Assert
				s.Error(err)
				s.Equal(expectedErr, err)
			}
		})
	}
}

// TestDeleteOrder_EmptyUUID verifies validation of empty UUID inputs.
// Tests that attempts to delete orders with empty UUIDs are properly rejected
// with ErrOrderNotFound, ensuring robust input validation.
func (s *OrdersServiceSuite) TestDeleteOrder_EmptyUUID() {
	// Arrange
	var (
		ctx         = context.Background()
		emptyUUID   = uuid.Nil
		expectedErr = model.ErrOrderNotFound
	)

	s.ordersRepository.On("DeleteOrder", ctx, emptyUUID).Return(expectedErr)
	// Act
	err := s.ordersService.DeleteOrder(ctx, emptyUUID)

	// Assert
	s.Error(err)
	s.Equal(err, expectedErr)
}
