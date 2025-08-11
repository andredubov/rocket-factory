package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/order/memory"
)

// TestDeleteOrder_Success verifies successful deletion of an existing order from the repository.
// Tests the complete flow: adding an order, deleting it, and verifying it can no longer be retrieved.
// Validates the repository's ability to permanently remove order records.
func TestDeleteOrder_Success(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()

		order = model.Order{
			OrderUUID: uuid.New(),
			Status:    model.OrderStatusPending,
		}
	)

	// Add order first
	err := ordersRepository.AddOrder(ctx, order)
	require.NoError(t, err)

	// Test
	err = ordersRepository.DeleteOrder(ctx, order.OrderUUID)

	// Verify
	require.NoError(t, err)

	// Verify order was actually deleted
	_, err = ordersRepository.GetOrder(ctx, order.OrderUUID)
	require.Equal(t, err, model.ErrOrderNotFoundWith(order.OrderUUID))
}

// TestDeleteOrder_NotFound verifies proper handling of delete attempts for non-existent orders.
// Tests that the repository returns ErrOrderNotFound when attempting to delete
// an order that doesn't exist, rather than silently succeeding.
func TestDeleteOrder_NotFound(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()
		nonExistentUUID  = uuid.New()
	)

	// Test
	err := ordersRepository.DeleteOrder(ctx, nonExistentUUID)

	// Verify
	require.Equal(t, err, model.ErrOrderNotFoundWith(nonExistentUUID))
}

// TestDeleteOrder_ConcurrentAccess verifies thread-safe deletion behavior under race conditions.
// Tests that concurrent delete operations for the same order result in exactly one
// successful deletion, maintaining data consistency.
func TestDeleteOrder_ConcurrentAccess(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()

		order = model.Order{
			OrderUUID: uuid.New(),
			Status:    model.OrderStatusPaid,
		}
	)

	// Add order first
	err := ordersRepository.AddOrder(ctx, order)
	require.NoError(t, err)

	// Test concurrent deletes
	var wg sync.WaitGroup
	wg.Add(2)

	var err1, err2 error
	go func() {
		defer wg.Done()
		err1 = ordersRepository.DeleteOrder(ctx, order.OrderUUID)
	}()
	go func() {
		defer wg.Done()
		err2 = ordersRepository.DeleteOrder(ctx, order.OrderUUID)
	}()

	wg.Wait()

	// Verify only one succeeded
	require.True(t, (err1 == nil && err2 != nil) || (err1 != nil && err2 == nil), "Exactly one DeleteOrder should succeed")

	// Verify order was deleted
	_, err = ordersRepository.GetOrder(ctx, order.OrderUUID)
	require.Equal(t, err, model.ErrOrderNotFoundWith(order.OrderUUID))
}
