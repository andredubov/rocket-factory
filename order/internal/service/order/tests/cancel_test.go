package tests

import (
	"context"

	"github.com/dvln/testify/mock"
	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// TestCancelOrder_Success tests successful order cancellation.
func (s *OrdersServiceSuite) TestCancelOrder_Success() {
	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
		order     = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPending,
		}
		expectedOrder = *order
	)
	expectedOrder.Status = model.OrderStatusCancelled

	// Mock expectations
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)
	s.ordersRepository.On("UpdateOrder", ctx, expectedOrder).Return(nil)

	// Test
	err := s.ordersService.CancelOrder(ctx, orderUUID)

	// Verify
	s.NoError(err)
	s.ordersRepository.AssertExpectations(s.T())
}

// TestCancelOrder_AlreadyPaid tests cancellation of already paid order.
func (s *OrdersServiceSuite) TestCancelOrder_AlreadyPaid() {
	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
		order     = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPaid,
		}
	)

	// Mock expectations
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)

	// Test
	err := s.ordersService.CancelOrder(ctx, orderUUID)

	// Verify
	s.Error(err)
	s.Equal(model.ErrOrderAlreadyPaid, err)
	s.ordersRepository.AssertExpectations(s.T())
	s.ordersRepository.AssertNotCalled(s.T(), "UpdateOrder")
}

// TestCancelOrder_AlreadyCancelled tests cancellation of already cancelled order
func (s *OrdersServiceSuite) TestCancelOrder_AlreadyCancelled() {
	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
		order     = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusCancelled,
		}
	)

	// Mock expectations
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)

	// Test
	err := s.ordersService.CancelOrder(ctx, orderUUID)

	// Verify
	s.Error(err)
	s.Equal(model.ErrOrderAlreadyCancelled, err)
	s.ordersRepository.AssertExpectations(s.T())
	s.ordersRepository.AssertNotCalled(s.T(), "UpdateOrder")
}

// TestCancelOrder_NotFound tests cancellation of non-existent order.
func (s *OrdersServiceSuite) TestCancelOrder_NotFound() {
	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
	)

	// Mock expectations
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	// Test
	err := s.ordersService.CancelOrder(ctx, orderUUID)

	// Verify
	s.Error(err)
	s.Equal(model.ErrOrderNotFound, err)
	s.ordersRepository.AssertExpectations(s.T())
	s.ordersRepository.AssertNotCalled(s.T(), "UpdateOrder")
}

// TestCancelOrder_UpdateError tests failure during order update.
func (s *OrdersServiceSuite) TestCancelOrder_UpdateError() {
	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
		order     = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPending,
		}
	)

	// Mock expectations
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)
	s.ordersRepository.On("UpdateOrder", ctx, mock.Anything).Return(model.ErrOrderNotFound)

	// Test
	err := s.ordersService.CancelOrder(ctx, orderUUID)

	// Verify
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.ordersRepository.AssertExpectations(s.T())
}
