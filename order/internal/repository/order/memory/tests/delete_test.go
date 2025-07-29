package tests

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository"
)

// TestDeleteOrder_Success verifies successful deletion of an existing order from the repository.
// Tests the complete flow: adding an order, deleting it, and verifying it can no longer be retrieved.
// Validates the repository's ability to permanently remove order records.
func (s *OrdersRepositorySuite) TestDeleteOrder_Success() {
	// Setup
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			Status:    model.OrderStatusPending,
		}
	)

	// Add order first
	err := s.ordersRepository.AddOrder(ctx, order)
	s.Require().NoError(err)

	// Test
	err = s.ordersRepository.DeleteOrder(ctx, order.OrderUUID)

	// Verify
	s.NoError(err)

	// Verify order was actually deleted
	_, err = s.ordersRepository.GetOrder(ctx, order.OrderUUID)
	s.Require().Equal(err, repository.ErrOrderNotFoundWith(order.OrderUUID))
}

// TestDeleteOrder_NotFound verifies proper handling of delete attempts for non-existent orders.
// Tests that the repository returns ErrOrderNotFound when attempting to delete
// an order that doesn't exist, rather than silently succeeding.
func (s *OrdersRepositorySuite) TestDeleteOrder_NotFound() {
	// Setup
	var (
		ctx             = context.Background()
		nonExistentUUID = uuid.New()
	)

	// Test
	err := s.ordersRepository.DeleteOrder(ctx, nonExistentUUID)

	// Verify
	s.Require().Equal(err, repository.ErrOrderNotFoundWith(nonExistentUUID))
}

// TestDeleteOrder_ConcurrentAccess verifies thread-safe deletion behavior under race conditions.
// Tests that concurrent delete operations for the same order result in exactly one
// successful deletion, maintaining data consistency.
func (s *OrdersRepositorySuite) TestDeleteOrder_ConcurrentAccess() {
	// Setup
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			Status:    model.OrderStatusPaid,
		}
	)

	// Add order first
	err := s.ordersRepository.AddOrder(ctx, order)
	s.Require().NoError(err)

	// Test concurrent deletes
	var wg sync.WaitGroup
	wg.Add(2)

	var err1, err2 error
	go func() {
		defer wg.Done()
		err1 = s.ordersRepository.DeleteOrder(ctx, order.OrderUUID)
	}()
	go func() {
		defer wg.Done()
		err2 = s.ordersRepository.DeleteOrder(ctx, order.OrderUUID)
	}()

	wg.Wait()

	// Verify only one succeeded
	s.True((err1 == nil && err2 != nil) || (err1 != nil && err2 == nil), "Exactly one DeleteOrder should succeed")

	// Verify order was deleted
	_, err = s.ordersRepository.GetOrder(ctx, order.OrderUUID)
	s.Require().Equal(err, repository.ErrOrderNotFoundWith(order.OrderUUID))
}
