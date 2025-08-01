package tests

import (
	"context"
	"testing"

	"github.com/dvln/testify/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/service/mocks"
	orders "github.com/andredubov/rocket-factory/order/internal/service/order"
)

// TestPayOrder_Success tests successful order payment scenario.
func TestPayOrder_Success(t *testing.T) {
	var (
		ordersRepository = mocks.NewOrdersRepository(t)
		paymentClient    = mocks.NewPaymentClient(t)
		inventoryClient  = mocks.NewInventoryClient(t)
		ordersService    = orders.NewService(ordersRepository, paymentClient, inventoryClient)
		ctx              = context.Background()
		orderUUID        = uuid.New()
		userUUID         = uuid.New()
		transactionUUID  = uuid.New()
		paymentMethod    = "CARD"
		order            = &model.Order{
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
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)
	paymentClient.On("PayOrder", ctx, order).Return(transactionUUID, nil)
	ordersRepository.On("UpdateOrder", ctx, expectedOrder).Return(nil)

	// Test
	result, err := ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, model.OrderStatusPaid, result.Status)
	require.Equal(t, transactionUUID, result.PaymentInfo.TransactionUUID)
	require.Equal(t, model.PaymentMethodCard, result.PaymentInfo.PaymentMethod)
	ordersRepository.AssertExpectations(t)
	paymentClient.AssertExpectations(t)
}

// TestPayOrder_InvalidStatus tests payment of order with invalid status.
func TestPayOrder_InvalidStatus(t *testing.T) {
	var (
		ordersRepository = mocks.NewOrdersRepository(t)
		paymentClient    = mocks.NewPaymentClient(t)
		inventoryClient  = mocks.NewInventoryClient(t)
		ordersService    = orders.NewService(ordersRepository, paymentClient, inventoryClient)
		ctx              = context.Background()
		orderUUID        = uuid.New()
		paymentMethod    = "CARD"

		order = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPaid, // Уже оплачен
		}
	)

	// Mock expectations
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)

	// Test
	result, err := ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, model.ErrInvalidOrderStatus, err)
	ordersRepository.AssertExpectations(t)
	paymentClient.AssertNotCalled(t, "PayOrder")
	ordersRepository.AssertNotCalled(t, "UpdateOrder")
}

// TestPayOrder_InvalidPaymentMethod tests payment with invalid method.
func TestPayOrder_InvalidPaymentMethod(t *testing.T) {
	var (
		ordersRepository = mocks.NewOrdersRepository(t)
		paymentClient    = mocks.NewPaymentClient(t)
		inventoryClient  = mocks.NewInventoryClient(t)
		ordersService    = orders.NewService(ordersRepository, paymentClient, inventoryClient)
		ctx              = context.Background()
		orderUUID        = uuid.New()
		paymentMethod    = "INVALID_METHOD"

		order = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPending,
		}
	)

	// Mock expectations
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)

	// Test
	result, err := ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, model.ErrInvalidPaymentMethod, err)
	ordersRepository.AssertExpectations(t)
	paymentClient.AssertNotCalled(t, "PayOrder")
	ordersRepository.AssertNotCalled(t, "UpdateOrder")
}

// TestPayOrder_PaymentFailed tests payment service failure scenario.
func TestPayOrder_PaymentFailed(t *testing.T) {
	var (
		ordersRepository = mocks.NewOrdersRepository(t)
		paymentClient    = mocks.NewPaymentClient(t)
		inventoryClient  = mocks.NewInventoryClient(t)
		ordersService    = orders.NewService(ordersRepository, paymentClient, inventoryClient)
		ctx              = context.Background()
		orderUUID        = uuid.New()
		paymentMethod    = "CARD"

		order = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPending,
		}
		paymentError = model.ErrOrderAlreadyPaid
	)

	// Mock expectations
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)
	paymentClient.On("PayOrder", ctx, order).Return(uuid.Nil, paymentError)

	// Test
	result, err := ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, paymentError, err)
	ordersRepository.AssertExpectations(t)
	paymentClient.AssertExpectations(t)
	ordersRepository.AssertNotCalled(t, "UpdateOrder")
}

// TestPayOrder_UpdateFailed tests order update failure after successful payment.
func TestPayOrder_UpdateFailed(t *testing.T) {
	var (
		ordersRepository = mocks.NewOrdersRepository(t)
		paymentClient    = mocks.NewPaymentClient(t)
		inventoryClient  = mocks.NewInventoryClient(t)
		ordersService    = orders.NewService(ordersRepository, paymentClient, inventoryClient)
		ctx              = context.Background()
		orderUUID        = uuid.New()
		userUUID         = uuid.New()
		transactionUUID  = uuid.New()
		paymentMethod    = "CARD"

		order = &model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPending,
		}
		updateError = model.ErrOrderNotFound
	)

	// Mock expectations
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)
	paymentClient.On("PayOrder", ctx, order).Return(transactionUUID, nil)
	ordersRepository.On("UpdateOrder", ctx, mock.Anything).Return(updateError)

	// Test
	result, err := ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, updateError, err)
	ordersRepository.AssertExpectations(t)
	paymentClient.AssertExpectations(t)
}

// TestPayOrder_GetOrderFromRepoFailed tests failure to retrieve order.
func TestPayOrder_GetOrderFromRepoFailed(t *testing.T) {
	var (
		ordersRepository = mocks.NewOrdersRepository(t)
		paymentClient    = mocks.NewPaymentClient(t)
		inventoryClient  = mocks.NewInventoryClient(t)
		ordersService    = orders.NewService(ordersRepository, paymentClient, inventoryClient)
		ctx              = context.Background()
		orderUUID        = uuid.New()
		paymentMethod    = "CARD"
		expectedError    = model.ErrOrderNotFound
	)

	// Mock expectations
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(nil, expectedError)

	// Test
	result, err := ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, expectedError, err)
	ordersRepository.AssertExpectations(t)
	paymentClient.AssertNotCalled(t, "PayOrder")
	ordersRepository.AssertNotCalled(t, "UpdateOrder")
}
