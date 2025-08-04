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

// TestCancelOrder_Success tests successful order cancellation.
func TestCancelOrder_Success(t *testing.T) {
	var (
		ordersRepository = mocks.NewOrdersRepository(t)
		paymentClient    = mocks.NewPaymentClient(t)
		inventoryClient  = mocks.NewInventoryClient(t)
		ordersService    = orders.NewService(ordersRepository, paymentClient, inventoryClient)
		ctx              = context.Background()
		orderUUID        = uuid.New()

		order = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPending,
		}
		expectedOrder = *order
	)
	expectedOrder.Status = model.OrderStatusCancelled

	// Mock expectations
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)
	ordersRepository.On("UpdateOrder", ctx, expectedOrder).Return(nil)

	// Test
	err := ordersService.CancelOrder(ctx, orderUUID)

	// Verify
	require.NoError(t, err)
	ordersRepository.AssertExpectations(t)
}

// TestCancelOrder_AlreadyPaid tests cancellation of already paid order.
func TestCancelOrder_AlreadyPaid(t *testing.T) {
	var (
		ordersRepository = mocks.NewOrdersRepository(t)
		paymentClient    = mocks.NewPaymentClient(t)
		inventoryClient  = mocks.NewInventoryClient(t)
		ordersService    = orders.NewService(ordersRepository, paymentClient, inventoryClient)
		ctx              = context.Background()
		orderUUID        = uuid.New()
		order            = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPaid,
		}
	)

	// Mock expectations
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)

	// Test
	err := ordersService.CancelOrder(ctx, orderUUID)

	// Verify
	require.Error(t, err)
	require.Equal(t, model.ErrOrderAlreadyPaid, err)
	ordersRepository.AssertExpectations(t)
	ordersRepository.AssertNotCalled(t, "UpdateOrder")
}

// TestCancelOrder_AlreadyCancelled tests cancellation of already cancelled order
func TestCancelOrder_AlreadyCancelled(t *testing.T) {
	var (
		ordersRepository = mocks.NewOrdersRepository(t)
		paymentClient    = mocks.NewPaymentClient(t)
		inventoryClient  = mocks.NewInventoryClient(t)
		ordersService    = orders.NewService(ordersRepository, paymentClient, inventoryClient)
		ctx              = context.Background()
		orderUUID        = uuid.New()
		order            = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusCancelled,
		}
	)

	// Mock expectations
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)

	// Test
	err := ordersService.CancelOrder(ctx, orderUUID)

	// Verify
	require.Error(t, err)
	require.Equal(t, model.ErrOrderAlreadyCancelled, err)
	ordersRepository.AssertExpectations(t)
	ordersRepository.AssertNotCalled(t, "UpdateOrder")
}

// TestCancelOrder_NotFound tests cancellation of non-existent order.
func TestCancelOrder_NotFound(t *testing.T) {
	var (
		ordersRepository = mocks.NewOrdersRepository(t)
		paymentClient    = mocks.NewPaymentClient(t)
		inventoryClient  = mocks.NewInventoryClient(t)
		ordersService    = orders.NewService(ordersRepository, paymentClient, inventoryClient)
		ctx              = context.Background()
		orderUUID        = uuid.New()
	)

	// Mock expectations
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	// Test
	err := ordersService.CancelOrder(ctx, orderUUID)

	// Verify
	require.Error(t, err)
	require.Equal(t, model.ErrOrderNotFound, err)
	ordersRepository.AssertExpectations(t)
	ordersRepository.AssertNotCalled(t, "UpdateOrder")
}

// TestCancelOrder_UpdateError tests failure during order update.
func TestCancelOrder_UpdateError(t *testing.T) {
	var (
		ordersRepository = mocks.NewOrdersRepository(t)
		paymentClient    = mocks.NewPaymentClient(t)
		inventoryClient  = mocks.NewInventoryClient(t)
		ordersService    = orders.NewService(ordersRepository, paymentClient, inventoryClient)
		ctx              = context.Background()
		orderUUID        = uuid.New()
		order            = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPending,
		}
	)

	// Mock expectations
	ordersRepository.On("GetOrder", ctx, orderUUID).Return(order, nil)
	ordersRepository.On("UpdateOrder", ctx, mock.Anything).Return(model.ErrOrderNotFound)

	// Test
	err := ordersService.CancelOrder(ctx, orderUUID)

	// Verify
	require.Error(t, err)
	require.ErrorIs(t, err, model.ErrOrderNotFound)
	ordersRepository.AssertExpectations(t)
}
