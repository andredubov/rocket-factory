package tests

import (
	"context"

	"github.com/dvln/testify/mock"
	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// TestPayOrder_Success tests successful order payment scenario.
func (s *OrdersServiceSuite) TestPayOrder_Success() {
	var (
		ctx             = context.Background()
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		transactionUUID = uuid.New()
		paymentMethod   = "CARD"
		order           = &model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPending,
		}
		expectedOrder = *order
	)

	// Подготовка ожидаемого состояния заказа
	expectedOrder.PaymentInfo = &model.PaymentInfo{
		PaymentMethod:   model.PaymentMethodCard,
		TransactionUUID: transactionUUID,
	}
	expectedOrder.Status = model.OrderStatusPaid

	// Mock expectations
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", ctx, order).Return(transactionUUID, nil)
	s.ordersRepository.On("UpdateOrder", ctx, expectedOrder).Return(nil)

	// Test
	result, err := s.ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	s.NoError(err)
	s.NotNil(result)
	s.Equal(model.OrderStatusPaid, result.Status)
	s.Equal(transactionUUID, result.PaymentInfo.TransactionUUID)
	s.Equal(model.PaymentMethodCard, result.PaymentInfo.PaymentMethod)
	s.ordersRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

// TestPayOrder_InvalidStatus tests payment of order with invalid status.
func (s *OrdersServiceSuite) TestPayOrder_InvalidStatus() {
	var (
		ctx           = context.Background()
		orderUUID     = uuid.New()
		paymentMethod = "CARD"
		order         = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPaid, // Уже оплачен
		}
	)

	// Mock expectations
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)

	// Test
	result, err := s.ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	s.Error(err)
	s.Nil(result)
	s.Equal(model.ErrInvalidOrderStatus, err)
	s.ordersRepository.AssertExpectations(s.T())
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder")
	s.ordersRepository.AssertNotCalled(s.T(), "UpdateOrder")
}

// TestPayOrder_InvalidPaymentMethod tests payment with invalid method.
func (s *OrdersServiceSuite) TestPayOrder_InvalidPaymentMethod() {
	var (
		ctx           = context.Background()
		orderUUID     = uuid.New()
		paymentMethod = "INVALID_METHOD"
		order         = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPending,
		}
	)

	// Mock expectations
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)

	// Test
	result, err := s.ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	s.Error(err)
	s.Nil(result)
	s.Equal(model.ErrInvalidPaymentMethod, err)
	s.ordersRepository.AssertExpectations(s.T())
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder")
	s.ordersRepository.AssertNotCalled(s.T(), "UpdateOrder")
}

// TestPayOrder_PaymentFailed tests payment service failure scenario.
func (s *OrdersServiceSuite) TestPayOrder_PaymentFailed() {
	var (
		ctx           = context.Background()
		orderUUID     = uuid.New()
		paymentMethod = "CARD"
		order         = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPending,
		}
		paymentError = model.ErrOrderAlreadyPaid
	)

	// Mock expectations
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", ctx, order).Return(uuid.Nil, paymentError)

	// Test
	result, err := s.ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	s.Error(err)
	s.Nil(result)
	s.Equal(paymentError, err)
	s.ordersRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
	s.ordersRepository.AssertNotCalled(s.T(), "UpdateOrder")
}

// TestPayOrder_UpdateFailed tests order update failure after successful payment.
func (s *OrdersServiceSuite) TestPayOrder_UpdateFailed() {
	var (
		ctx             = context.Background()
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		transactionUUID = uuid.New()
		paymentMethod   = "CARD"
		order           = &model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPending,
		}
		updateError = model.ErrOrderNotFound
	)

	// Mock expectations
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", ctx, order).Return(transactionUUID, nil)
	s.ordersRepository.On("UpdateOrder", ctx, mock.Anything).Return(updateError)

	// Test
	result, err := s.ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	s.Error(err)
	s.Nil(result)
	s.Equal(updateError, err)
	s.ordersRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

// TestPayOrder_GetOrderFromRepoFailed tests failure to retrieve order.
func (s *OrdersServiceSuite) TestPayOrder_GetOrderFromRepoFailed() {
	var (
		ctx           = context.Background()
		orderUUID     = uuid.New()
		paymentMethod = "CARD"
		expectedError = model.ErrOrderNotFound
	)

	// Mock expectations
	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(nil, expectedError)

	// Test
	result, err := s.ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedError, err)
	s.ordersRepository.AssertExpectations(s.T())
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder")
	s.ordersRepository.AssertNotCalled(s.T(), "UpdateOrder")
}
