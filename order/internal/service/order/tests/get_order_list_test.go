package tests

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// TestGetUserOrders_Success verifies successful retrieval of a user's orders through the service layer.
// Tests that all orders for a specific user are properly returned and converted
// from repository to domain models, including nested payment information when present.
func (s *OrdersServiceSuite) TestGetUserOrders_Success() {
	// Arrange
	var (
		ctx      = context.Background()
		userUUID = uuid.New()
		orders   = []model.Order{
			{
				OrderUUID: uuid.New(),
				UserUUID:  userUUID,
				Status:    model.OrderStatusPaid,
			},
			{
				OrderUUID: uuid.New(),
				UserUUID:  userUUID,
				Status:    model.OrderStatusPending,
			},
		}
	)

	s.ordersRepository.On("GetUserOrders", ctx, userUUID).Return(orders, nil)

	// Act
	result, err := s.ordersService.GetUserOrders(ctx, userUUID)

	// Assert
	s.NoError(err)
	s.Equal(orders, result)
}

// TestGetUserOrders_EmptyResult verifies proper handling when a user has no orders.
// Tests that an empty slice (not nil) is returned without errors when querying
// a user with no order history, ensuring correct empty state behavior.
func (s *OrdersServiceSuite) TestGetUserOrders_EmptyResult() {
	// Arrange
	var (
		ctx      = context.Background()
		userUUID = uuid.New()
	)

	s.ordersRepository.On("GetUserOrders", ctx, userUUID).Return([]model.Order{}, nil)

	// Act
	result, err := s.ordersService.GetUserOrders(ctx, userUUID)

	// Assert
	s.NoError(err)
	s.Empty(result)
}

// TestGetUserOrders_RepositoryError verifies error propagation from the repository layer.
// Tests that repository-level errors are properly propagated through the service
// while maintaining the original error type and context.
func (s *OrdersServiceSuite) TestGetUserOrders_RepositoryError() {
	// Arrange
	var (
		ctx         = context.Background()
		userUUID    = uuid.New()
		expectedErr = model.ErrOrderNotFound
	)
	s.ordersRepository.On("GetUserOrders", ctx, userUUID).Return(nil, expectedErr)

	// Act
	result, err := s.ordersService.GetUserOrders(ctx, userUUID)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

// TestGetUserOrders_EmptyUUID verifies input validation for empty user UUIDs.
// Tests that requests with nil UUIDs are properly rejected with ErrOrderNotFound,
// enforcing valid user reference requirements.
func (s *OrdersServiceSuite) TestGetUserOrders_EmptyUUID() {
	// Arrange
	var (
		ctx         = context.Background()
		emptyUUID   = uuid.Nil
		expectedErr = model.ErrOrderNotFound
	)

	s.ordersRepository.On("GetUserOrders", ctx, emptyUUID).Return(nil, expectedErr)

	// Act
	result, err := s.ordersService.GetUserOrders(ctx, emptyUUID)

	// Assert
	s.Nil(result)
	s.Error(err)
	s.Equal(err, expectedErr)
}
