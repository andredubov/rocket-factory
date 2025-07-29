package tests

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository"
)

// TestAddOrder_Success verifies that a valid order can be successfully added to the repository.
// Tests the happy path scenario where all required order fields are properly populated
// and verifies the order can be retrieved after creation.
func (s *OrdersRepositorySuite) TestAddOrder_Success() {
	// Setup
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID:  uuid.New(),
			UserUUID:   uuid.New(),
			PartUUIDs:  []uuid.UUID{uuid.New()},
			TotalPrice: 100.50,
			Status:     model.OrderStatusPending,
		}
	)

	// Test
	err := s.ordersRepository.AddOrder(ctx, order)

	// Verify
	s.NoError(err)

	// Verify order was actually added
	retrieved, err := s.ordersRepository.GetOrder(ctx, order.OrderUUID)
	s.Require().NoError(err)
	s.Require().Equal(order.OrderUUID, retrieved.OrderUUID)
	s.Require().Equal(order.Status, retrieved.Status)
}

// TestAddOrder_DuplicateOrder verifies the repository correctly handles duplicate order attempts.
// Tests that attempting to add an order with an existing UUID returns the expected
// ErrOrderAlreadyExists error, maintaining data integrity.
func (s *OrdersRepositorySuite) TestAddOrder_DuplicateOrder() {
	// Setup
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			Status:    model.OrderStatusPaid,
		}
	)

	// Add order first time
	err := s.ordersRepository.AddOrder(ctx, order)
	s.Require().NoError(err)

	// Test adding duplicate
	err = s.ordersRepository.AddOrder(ctx, order)

	// Verify
	s.Require().Equal(err, repository.ErrOrderAlreadyExistsWith(order.OrderUUID))
}

// TestAddOrder_ConcurrentAccess verifies thread-safe order creation behavior.
// Tests that concurrent attempts to create the same order result in exactly one
// successful creation, ensuring data consistency under race conditions.
func (s *OrdersRepositorySuite) TestAddOrder_ConcurrentAccess() {
	// Setup
	var (
		ctx   = context.Background()
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
		err1 = s.ordersRepository.AddOrder(ctx, order)
	}()
	go func() {
		defer wg.Done()
		err2 = s.ordersRepository.AddOrder(ctx, order)
	}()

	wg.Wait()

	// Verify only one succeeded
	s.True((err1 == nil && err2 != nil) || (err1 != nil && err2 == nil), "Exactly one AddOrder should succeed")

	// Verify order exists
	_, err := s.ordersRepository.GetOrder(ctx, order.OrderUUID)
	s.Require().NoError(err)
}
