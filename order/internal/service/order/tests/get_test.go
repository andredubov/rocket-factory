package tests

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// TestGetOrder_Success tests successful retrieval of an existing order.
func (s *OrdersServiceSuite) TestGetOrder_Success() {
	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
		order     = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPending,
		}
	)

	// Mock expectation
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)

	// Test
	result, err := s.ordersService.GetOrder(ctx, orderUUID)

	// Verify
	s.NoError(err)
	s.Equal(order, result)
	s.ordersRepository.AssertExpectations(s.T())
}

// TestGetOrder_NotFound tests retrieval of a non-existent order.
func (s *OrdersServiceSuite) TestGetOrder_NotFound() {
	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
	)

	// Mock expectation
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	// Test
	result, err := s.ordersService.GetOrder(ctx, orderUUID)

	// Verify
	s.Error(err)
	s.Nil(result)
	s.Equal(model.ErrOrderNotFound, err)
	s.ordersRepository.AssertExpectations(s.T())
}

// TestGetOrder_RepositoryError tests error handling during order retrieval.
func (s *OrdersServiceSuite) TestGetOrder_RepositoryError() {
	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
		repoError = model.ErrOrderNotFound
	)

	// Mock expectation
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(nil, repoError)

	// Test
	result, err := s.ordersService.GetOrder(ctx, orderUUID)

	// Verify
	s.Error(err)
	s.Nil(result)
	s.Equal(repoError, err)
	s.ordersRepository.AssertExpectations(s.T())
}
