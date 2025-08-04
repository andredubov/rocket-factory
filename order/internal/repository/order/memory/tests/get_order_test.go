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

// TestGetUserOrders_Success verifies successful retrieval of a user's orders from the repository.
// Tests that all orders for a specific user are returned while excluding other users' orders.
// Validates correct filtering by UserUUID and proper order data structure.
func TestGetUserOrders_Success(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()
		user1            = uuid.New()
		user2            = uuid.New()

		// Create test orders
		orders = []model.Order{
			{
				OrderUUID: uuid.New(),
				UserUUID:  user1,
				Status:    model.OrderStatusPending,
			},
			{
				OrderUUID: uuid.New(),
				UserUUID:  user1,
				Status:    model.OrderStatusPaid,
			},
			{
				OrderUUID: uuid.New(),
				UserUUID:  user2,
				Status:    model.OrderStatusCancelled,
			},
		}
	)

	// Add all orders to repository
	for _, order := range orders {
		err := ordersRepository.AddOrder(ctx, order)
		require.NoError(t, err)
	}

	// Test
	result, err := ordersRepository.GetUserOrders(ctx, user1)

	// Verify
	require.NoError(t, err)
	require.Len(t, result, 2)

	for _, order := range result {
		require.Equal(t, user1, order.UserUUID)
	}
}

// TestGetUserOrders_NoOrders verifies correct handling when a user has no orders.
// Tests that the repository returns an empty slice (not nil) when no orders exist
// for the requested user, without returning an error.
func TestGetUserOrders_NoOrders(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()
		userWithNoOrders = uuid.New()
	)

	// Test
	result, err := ordersRepository.GetUserOrders(ctx, userWithNoOrders)

	// Verify
	require.NoError(t, err)
	require.Empty(t, result)
}

// TestGetUserOrders_ConcurrentAccess verifies thread-safe order retrieval under concurrent access.
// Tests that the repository correctly handles simultaneous read operations while
// maintaining data consistency across goroutines.
func TestGetUserOrders_ConcurrentAccess(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()
		user             = uuid.New()
	)

	// Add orders in separate goroutines
	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			order := model.Order{
				OrderUUID: uuid.New(),
				UserUUID:  user,
				Status:    model.OrderStatusPending,
			}
			err := ordersRepository.AddOrder(ctx, order)
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	// Test concurrent reads
	wg.Add(2)
	var result1, result2 []model.Order
	var err1, err2 error

	go func() {
		defer wg.Done()
		result1, err1 = ordersRepository.GetUserOrders(ctx, user)
	}()
	go func() {
		defer wg.Done()
		result2, err2 = ordersRepository.GetUserOrders(ctx, user)
	}()

	wg.Wait()

	// Verify
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.Len(t, result1, 3)
	require.Len(t, result2, 3)
}

// TestGetUserOrders_ReturnsCopies verifies the repository returns defensive copies of orders.
// Tests that modifications to returned order objects don't affect the stored data,
// ensuring data integrity.
func TestGetUserOrders_ReturnsCopies(t *testing.T) {
	// Setup
	var (
		ordersRepository = memory.NewOrderRepository()
		ctx              = context.Background()
		user             = uuid.New()
		order            = model.Order{
			OrderUUID: uuid.New(),
			UserUUID:  user,
			Status:    model.OrderStatusPending,
		}
	)

	err := ordersRepository.AddOrder(ctx, order)
	require.NoError(t, err)

	// Test
	result, err := ordersRepository.GetUserOrders(ctx, user)
	require.NoError(t, err)
	require.Len(t, result, 1)

	// Modify the returned order
	result[0].Status = model.OrderStatusPaid

	// Verify original wasn't modified
	storedOrder, err := ordersRepository.GetOrder(ctx, order.OrderUUID)
	require.NoError(t, err)
	require.Equal(t, model.OrderStatusPending, storedOrder.Status)
}
