package tests

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// TestGetOrder_Success verifies successful retrieval of an order through the service layer.
// Tests that a valid order UUID returns the corresponding order with all fields properly
// converted from the repository model to the domain model.
func (s *OrdersServiceSuite) TestGetOrder_Success() {
	// Arrange
	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
		order     = &model.Order{
			OrderUUID: orderUUID,
			UserUUID:  uuid.New(),
			Status:    model.OrderStatusPaid,
		}
	)

	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)

	// Act
	retrivedOrder, err := s.ordersService.GetOrder(ctx, orderUUID)

	// Assert
	s.NoError(err)
	s.Equal(retrivedOrder, order)
}

// TestGetOrder_NotFound verifies proper handling of requests for non-existent orders.
// Tests that the service returns ErrOrderNotFound when attempting to retrieve
// an order that doesn't exist, maintaining proper error handling semantics.
func (s *OrdersServiceSuite) TestGetOrder_NotFound() {
	// Arrange
	var (
		ctx         = context.Background()
		orderUUID   = uuid.New()
		expectedErr = model.ErrOrderNotFound
	)

	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(nil, expectedErr)

	// Act
	result, err := s.ordersService.GetOrder(ctx, orderUUID)

	// Assert
	s.Nil(result)
	s.Error(err)
	s.Equal(err, expectedErr)
}

// TestGetOrder_EmptyUUID verifies validation of empty UUID inputs.
// Tests that attempts to retrieve orders with empty UUIDs are properly rejected
// with ErrOrderNotFound, ensuring robust input validation.
func (s *OrdersServiceSuite) TestGetOrder_EmptyUUID() {
	// Arrange
	var (
		ctx         = context.Background()
		emptyUUID   = uuid.Nil
		expectedErr = model.ErrOrderNotFound
	)

	s.ordersRepository.On("GetOrder", ctx, emptyUUID).Return(nil, expectedErr)

	// Act
	result, err := s.ordersService.GetOrder(ctx, emptyUUID)

	// Assert
	s.Nil(result)
	s.Error(err)
	s.Equal(err, expectedErr)
}
