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

// TestAddOrder_Success verifies that a valid order can be successfully added to the repository.
// Tests the happy path scenario where all required order fields are properly populated
// and verifies the order can be retrieved after creation.
func TestAddOrder_Success(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()
		order            = model.Order{
			OrderUUID:  uuid.New(),
			UserUUID:   uuid.New(),
			PartUUIDs:  []uuid.UUID{uuid.New()},
			TotalPrice: 100.50,
			Status:     model.OrderStatusPending,
		}
	)

	// Test
	err := ordersRepository.AddOrder(ctx, order)

	// Verify
	require.NoError(t, err)

	// Verify order was actually added
	retrieved, err := ordersRepository.GetOrder(ctx, order.OrderUUID)
	require.NoError(t, err)
	require.Equal(t, order.OrderUUID, retrieved.OrderUUID)
	require.Equal(t, order.Status, retrieved.Status)
}

// TestAddOrder_DuplicateOrder verifies the repository correctly handles duplicate order attempts.
// Tests that attempting to add an order with an existing UUID returns the expected
// ErrOrderAlreadyExists error, maintaining data integrity.
func TestAddOrder_DuplicateOrder(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()
		order            = model.Order{
			OrderUUID: uuid.New(),
			Status:    model.OrderStatusPaid,
		}
	)

	// Add order first time
	err := ordersRepository.AddOrder(ctx, order)
	require.NoError(t, err)

	// Test adding duplicate
	err = ordersRepository.AddOrder(ctx, order)

	// Verify
	require.Equal(t, err, model.ErrOrderAlreadyExistsWith(order.OrderUUID))
}

// TestAddOrder_ConcurrentAccess verifies thread-safe order creation behavior.
// Tests that concurrent attempts to create the same order result in exactly one
// successful creation, ensuring data consistency under race conditions.
func TestAddOrder_ConcurrentAccess(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()

		order = model.Order{
			OrderUUID: uuid.New(),
			Status:    model.OrderStatusPending,
		}
	)

	// Test concurrent adds
	var wg sync.WaitGroup
	wg.Add(2)

	var err1, err2 error
	go func() {
		defer wg.Done()
		err1 = ordersRepository.AddOrder(ctx, order)
	}()
	go func() {
		defer wg.Done()
		err2 = ordersRepository.AddOrder(ctx, order)
	}()

	wg.Wait()

	// Verify only one succeeded
	require.True(t, (err1 == nil && err2 != nil) || (err1 != nil && err2 == nil), "Exactly one AddOrder should succeed")

	// Verify order exists
	_, err := ordersRepository.GetOrder(ctx, order.OrderUUID)
	require.NoError(t, err)
}
