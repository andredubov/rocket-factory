package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/order/memory"
)

// TestUpdateOrder_Success verifies successful order updates in the repository.
// Tests the complete flow: adding an order, updating its status and payment info,
// and verifying the changes persist. Validates field-level updates includin status transitions and payment information.
func TestUpdateOrder_Success(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()
		originalOrder    = model.Order{
			OrderUUID:  uuid.New(),
			UserUUID:   uuid.New(),
			Status:     model.OrderStatusPending,
			TotalPrice: 100.0,
		}
		newTotalPrice = 150.0
	)

	// Add initial order
	err := ordersRepository.AddOrder(ctx, originalOrder)
	require.NoError(t, err)

	updatedOrder := model.OrderUpdateInfo{
		OrderUUID:  originalOrder.OrderUUID,
		UserUUID:   &originalOrder.UserUUID,
		Status:     lo.ToPtr(model.OrderStatusPaid),
		TotalPrice: &newTotalPrice,
		PaymentInfo: &model.PaymentInfo{
			PaymentMethod:   model.PaymentMethodCard,
			TransactionUUID: uuid.New(),
		},
	}

	// Test
	err = ordersRepository.UpdateOrder(ctx, updatedOrder)

	// Verify
	require.NoError(t, err)

	// Verify changes were applied
	retrieved, err := ordersRepository.GetOrder(ctx, originalOrder.OrderUUID)
	require.NoError(t, err)
	require.Equal(t, *updatedOrder.Status, retrieved.Status)
	require.Equal(t, *updatedOrder.TotalPrice, retrieved.TotalPrice)
	require.Equal(t, updatedOrder.PaymentInfo.PaymentMethod, retrieved.PaymentInfo.PaymentMethod)
}

// TestUpdateOrder_NotFound verifies proper handling of update attempts for non-existent orders.
// Tests that the repository returns ErrOrderNotFound when attempting to update
// an order that doesn't exist, ensuring data integrity.
func TestUpdateOrder_NotFound(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()
		nonExistentOrder = model.OrderUpdateInfo{
			OrderUUID: uuid.New(),
			Status:    lo.ToPtr(model.OrderStatusPending),
		}
	)

	// Test
	err := ordersRepository.UpdateOrder(ctx, nonExistentOrder)

	// Verify
	require.Equal(t, err, model.ErrOrderNotFoundWith(nonExistentOrder.OrderUUID))
}

// TestUpdateOrder_ConcurrentAccess verifies thread-safe update behavior under race conditions.
// Tests that concurrent updates to the same order are handled correctly, with
// the repository applying one of the updates while maintaining data consistency.
func TestUpdateOrder_ConcurrentAccess(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()

		originalOrder = model.Order{
			OrderUUID:   uuid.New(),
			Status:      model.OrderStatusPending,
			PaymentInfo: &model.PaymentInfo{PaymentMethod: model.PaymentMethodCard},
		}
		newTotalPrice1 = 300.0
		newTotalPrice2 = 200.0
	)

	// Add initial order
	err := ordersRepository.AddOrder(ctx, originalOrder)
	require.NoError(t, err)

	// Prepare updates
	update1 := model.OrderUpdateInfo{
		OrderUUID: originalOrder.OrderUUID,
		UserUUID:  &originalOrder.UserUUID,
		Status:    &originalOrder.Status,
		PaymentInfo: &model.PaymentInfo{
			PaymentMethod:   model.PaymentMethodCard,
			TransactionUUID: uuid.New(),
		},
		TotalPrice: &newTotalPrice1,
	}

	update2 := model.OrderUpdateInfo{
		OrderUUID: originalOrder.OrderUUID,
		UserUUID:  &originalOrder.UserUUID,
		Status:    &originalOrder.Status,
		PaymentInfo: &model.PaymentInfo{
			PaymentMethod:   model.PaymentMethodCard,
			TransactionUUID: uuid.New(),
		},
		TotalPrice: &newTotalPrice2,
	}

	// Test concurrent updates
	var wg sync.WaitGroup
	wg.Add(2)

	var err1, err2 error
	go func() {
		defer wg.Done()
		err1 = ordersRepository.UpdateOrder(ctx, update1)
	}()
	go func() {
		defer wg.Done()
		err2 = ordersRepository.UpdateOrder(ctx, update2)
	}()

	wg.Wait()

	// Verify both succeeded
	require.True(t, (err1 == nil && err2 == nil))

	// Verify order was updated
	retrieved, err := ordersRepository.GetOrder(ctx, originalOrder.OrderUUID)
	require.NoError(t, err)

	// Check which update was applied
	if err1 == nil {
		require.Equal(t, *update1.Status, retrieved.Status)
	} else {
		require.Equal(t, *update2.TotalPrice, retrieved.TotalPrice)
	}
}
