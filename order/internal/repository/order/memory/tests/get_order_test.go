package tests

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// TestGetUserOrders_Success verifies successful retrieval of a user's orders from the repository.
// Tests that all orders for a specific user are returned while excluding other users' orders.
// Validates correct filtering by UserUUID and proper order data structure.
func (s *OrdersRepositorySuite) TestGetUserOrders_Success() {
	// Setup
	var (
		ctx   = context.Background()
		user1 = uuid.New()
		user2 = uuid.New()

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
		err := s.ordersRepository.AddOrder(ctx, order)
		s.Require().NoError(err)
	}

	// Test
	result, err := s.ordersRepository.GetUserOrders(ctx, user1)

	// Verify
	s.Require().NoError(err)
	s.Require().Len(result, 2)

	for _, order := range result {
		s.Require().Equal(user1, order.UserUUID)
	}
}

// TestGetUserOrders_NoOrders verifies correct handling when a user has no orders.
// Tests that the repository returns an empty slice (not nil) when no orders exist
// for the requested user, without returning an error.
func (s *OrdersRepositorySuite) TestGetUserOrders_NoOrders() {
	// Setup
	var (
		ctx              = context.Background()
		userWithNoOrders = uuid.New()
	)

	// Test
	result, err := s.ordersRepository.GetUserOrders(ctx, userWithNoOrders)

	// Verify
	s.Require().NoError(err)
	s.Require().Empty(result)
}

// TestGetUserOrders_ConcurrentAccess verifies thread-safe order retrieval under concurrent access.
// Tests that the repository correctly handles simultaneous read operations while
// maintaining data consistency across goroutines.
func (s *OrdersRepositorySuite) TestGetUserOrders_ConcurrentAccess() {
	// Setup
	var (
		ctx  = context.Background()
		user = uuid.New()
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
			err := s.ordersRepository.AddOrder(ctx, order)
			s.NoError(err)
		}()
	}
	wg.Wait()

	// Test concurrent reads
	wg.Add(2)
	var result1, result2 []model.Order
	var err1, err2 error

	go func() {
		defer wg.Done()
		result1, err1 = s.ordersRepository.GetUserOrders(ctx, user)
	}()
	go func() {
		defer wg.Done()
		result2, err2 = s.ordersRepository.GetUserOrders(ctx, user)
	}()

	wg.Wait()

	// Verify
	s.Require().NoError(err1)
	s.Require().NoError(err2)
	s.Require().Len(result1, 3)
	s.Require().Len(result2, 3)
}

// TestGetUserOrders_ReturnsCopies verifies the repository returns defensive copies of orders.
// Tests that modifications to returned order objects don't affect the stored data,
// ensuring data integrity.
func (s *OrdersRepositorySuite) TestGetUserOrders_ReturnsCopies() {
	// Setup
	var (
		ctx   = context.Background()
		user  = uuid.New()
		order = model.Order{
			OrderUUID: uuid.New(),
			UserUUID:  user,
			Status:    model.OrderStatusPending,
		}
	)

	err := s.ordersRepository.AddOrder(ctx, order)
	s.Require().NoError(err)

	// Test
	result, err := s.ordersRepository.GetUserOrders(ctx, user)
	s.Require().NoError(err)
	s.Require().Len(result, 1)

	// Modify the returned order
	result[0].Status = model.OrderStatusPaid

	// Verify original wasn't modified
	storedOrder, err := s.ordersRepository.GetOrder(ctx, order.OrderUUID)
	s.Require().NoError(err)
	s.Require().Equal(model.OrderStatusPending, storedOrder.Status)
}
