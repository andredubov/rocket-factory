package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/service/mocks"
	orders "github.com/andredubov/rocket-factory/order/internal/service/order"
)

// TestGetOrder_Success tests successful retrieval of an existing order.
func TestGetOrder_Success(t *testing.T) {
	var (
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		orderUUID              = uuid.New()
		order                  = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPending,
		}
	)

	// Mock expectation
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)

	// Test
	result, err := ordersService.GetOrder(ctx, orderUUID)

	// Verify
	require.NoError(t, err)
	require.Equal(t, order, result)
	ordersRepository.AssertExpectations(t)
}

// TestGetOrder_NotFound tests retrieval of a non-existent order.
func TestGetOrder_NotFound(t *testing.T) {
	var (
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		orderUUID              = uuid.New()
	)

	// Mock expectation
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	// Test
	result, err := ordersService.GetOrder(ctx, orderUUID)

	// Verify
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, model.ErrOrderNotFound, err)
	ordersRepository.AssertExpectations(t)
}

// TestGetOrder_RepositoryError tests error handling during order retrieval.
func TestGetOrder_RepositoryError(t *testing.T) {
	var (
		ordersRepository       = mocks.NewOrdersRepository(t)
		paymentClient          = mocks.NewPaymentClient(t)
		inventoryClient        = mocks.NewInventoryClient(t)
		orderPaidEventProducer = mocks.NewProducerService(t)
		ordersService          = orders.NewService(ordersRepository, paymentClient, inventoryClient, orderPaidEventProducer)
		ctx                    = context.Background()
		orderUUID              = uuid.New()
		repoError              = model.ErrOrderNotFound
	)

	// Mock expectation
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(nil, repoError)

	// Test
	result, err := ordersService.GetOrder(ctx, orderUUID)

	// Verify
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, repoError, err)
	ordersRepository.AssertExpectations(t)
}
