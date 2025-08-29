package tests

import (
	"context"
	"errors"
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
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		orderUUID              = uuid.New()
		userUUID               = uuid.New()
		transactionUUID        = uuid.New()
		paymentMethod          = "CARD"
		order                  = &model.Order{
			OrderUUID:   orderUUID,
			UserUUID:    userUUID,
			Status:      model.OrderStatusPending,
			PaymentInfo: nil,
			PartUUIDs:   nil,
			TotalPrice:  10.0,
		}
		expectedOrder = *order
	)

	// Подготовка ожидаемого состояния заказа
	expectedOrder.PaymentInfo = &model.PaymentInfo{
		PaymentMethod:   model.PaymentMethodCard,
		TransactionUUID: transactionUUID,
	}
	expectedOrder.Status = model.OrderStatusPaid

	updateInfo := model.OrderUpdateInfo{
		OrderUUID:  expectedOrder.OrderUUID,
		UserUUID:   &expectedOrder.UserUUID,
		TotalPrice: &expectedOrder.TotalPrice,
		Status:     &expectedOrder.Status,
		PartUUIDs:  nil,
		PaymentInfo: &model.PaymentInfo{
			PaymentMethod:   model.PaymentMethodCard,
			TransactionUUID: transactionUUID,
		},
	}

	event := model.OrderPaidEvent{
		UUID:            expectedOrder.OrderUUID.String(),
		OrderUUID:       expectedOrder.OrderUUID.String(),
		UserUUID:        expectedOrder.UserUUID.String(),
		PaymentMethod:   expectedOrder.PaymentInfo.PaymentMethod,
		TransactionUUID: expectedOrder.PaymentInfo.TransactionUUID.String(),
	}

	// Mock expectations
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)
	paymentClient.On("PayOrder", ctx, order).Return(transactionUUID, nil)
	ordersRepository.On("UpdateOrder", ctx, updateInfo).Return(nil)
	orderPaidEventProducer.On("ProduceOrderPaidEvent", ctx, event).Return(nil)

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
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		orderUUID              = uuid.New()
		paymentMethod          = "CARD"

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
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		orderUUID              = uuid.New()
		paymentMethod          = "INVALID_METHOD"

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
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		orderUUID              = uuid.New()
		paymentMethod          = "CARD"

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
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		orderUUID              = uuid.New()
		userUUID               = uuid.New()
		transactionUUID        = uuid.New()
		paymentMethod          = "CARD"

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
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		orderUUID              = uuid.New()
		paymentMethod          = "CARD"
		expectedError          = model.ErrOrderNotFound
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

// TestPayOrder_ProduceEventFailed tests failure to produce order paid event.
func TestPayOrder_ProduceEventFailed(t *testing.T) {
	var (
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		orderUUID              = uuid.New()
		userUUID               = uuid.New()
		transactionUUID        = uuid.New()
		paymentMethod          = "CARD"

		order = &model.Order{
			OrderUUID:   orderUUID,
			UserUUID:    userUUID,
			Status:      model.OrderStatusPending,
			PaymentInfo: nil,
			PartUUIDs:   nil,
			TotalPrice:  10.0,
		}
		produceError = errors.New("kafka producer error")
	)

	// Подготовка ожидаемого события
	expectedEvent := model.OrderPaidEvent{
		UUID:            orderUUID.String(),
		OrderUUID:       orderUUID.String(),
		UserUUID:        userUUID.String(),
		PaymentMethod:   model.PaymentMethodCard,
		TransactionUUID: transactionUUID.String(),
	}

	// Mock expectations
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)
	paymentClient.On("PayOrder", ctx, order).Return(transactionUUID, nil)
	ordersRepository.On("UpdateOrder", ctx, mock.Anything).Return(nil)
	orderPaidEventProducer.On("ProduceOrderPaidEvent", ctx, expectedEvent).Return(produceError)

	// Test
	result, err := ordersService.PayOrder(ctx, orderUUID, paymentMethod)

	// Verify
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, produceError, err)
	ordersRepository.AssertExpectations(t)
	paymentClient.AssertExpectations(t)
	orderPaidEventProducer.AssertExpectations(t)
}
