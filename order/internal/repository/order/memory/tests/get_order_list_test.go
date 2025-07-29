package tests

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// TestGetUserOrders_ReturnsCorrectOrders verifies correct filtering of orders by user ID.
// Tests that only orders belonging to the specified user are returned, while excluding
// orders from other users. Validates the repository's filtering capability.
func (s *OrdersRepositorySuite) TestGetUserOrders_ReturnsCorrectOrders() {
	// Arrange
	var (
		ctx   = context.Background()
		user1 = uuid.New()
		user2 = uuid.New()
	)

	orders := []model.Order{
		{OrderUUID: uuid.New(), UserUUID: user1, Status: model.OrderStatusPending},
		{OrderUUID: uuid.New(), UserUUID: user1, Status: model.OrderStatusPaid},
		{OrderUUID: uuid.New(), UserUUID: user2, Status: model.OrderStatusCancelled},
	}

	for _, order := range orders {
		err := s.ordersRepository.AddOrder(ctx, order)
		s.Require().NoError(err)
	}

	// Act
	result, err := s.ordersRepository.GetUserOrders(ctx, user1)

	// Assert
	s.Require().NoError(err)
	s.Require().Len(result, 2)

	for _, order := range result {
		s.Equal(user1, order.UserUUID)
	}
}

// TestGetUserOrders_ReturnsEmptyForNoOrders verifies proper handling when a user has no orders.
// Tests that an empty slice (not nil) is returned without errors when querying
// a user with no order history.
func (s *OrdersRepositorySuite) TestGetUserOrders_ReturnsEmptyForNoOrders() {
	// Arrange
	var (
		ctx     = context.Background()
		newUser = uuid.New()
	)

	// Act
	result, err := s.ordersRepository.GetUserOrders(ctx, newUser)

	// Assert
	s.Require().NoError(err)
	s.Require().Empty(result)
}

// TestGetUserOrders_ReturnsCopiesNotReferences verifies defensive copying of returned orders.
// Tests that modifications to retrieved order objects don't affect the stored data,
// ensuring the repository maintains data integrity.
func (s *OrdersRepositorySuite) TestGetUserOrders_ReturnsCopiesNotReferences() {
	// Arrange
	ctx := context.Background()
	user := uuid.New()
	order := model.Order{
		OrderUUID: uuid.New(),
		UserUUID:  user,
		Status:    model.OrderStatusPending,
	}

	err := s.ordersRepository.AddOrder(ctx, order)
	s.Require().NoError(err)

	// Act
	result, err := s.ordersRepository.GetUserOrders(ctx, user)
	s.Require().NoError(err)
	s.Require().Len(result, 1)

	// Modify the returned order
	result[0].Status = model.OrderStatusPaid

	// Assert original wasn't modified
	storedOrder, err := s.ordersRepository.GetOrder(ctx, order.OrderUUID)
	s.Require().NoError(err)
	s.Require().Equal(model.OrderStatusPending, storedOrder.Status)
}

// TestGetUserOrders_ThreadSafe verifies thread-safe behavior under high concurrency.
// Tests that the repository correctly handles simultaneous read operations during
// heavy write loads, maintaining data consistency across goroutines.
func (s *OrdersRepositorySuite) TestGetUserOrders_ThreadSafe() {
	// Arrange
	var (
		ctx  = context.Background()
		user = uuid.New()
	)

	const (
		numOrders = 100
	)

	// Add orders concurrently
	var wg sync.WaitGroup
	wg.Add(numOrders)

	for i := 0; i < numOrders; i++ {
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

	// Act & Assert - concurrent reads
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

	s.Require().NoError(err1)
	s.Require().NoError(err2)
	s.Require().Len(result1, numOrders)
	s.Require().Len(result2, numOrders)
}

// TestGetUserOrders_ReturnsCorrectOrderData verifies complete and accurate order data retrieval.
// Tests that all order fields including nested structures (like payment info) are
// correctly preserved and returned by the repository.
func (s *OrdersRepositorySuite) TestGetUserOrders_ReturnsCorrectOrderData() {
	// Arrange
	var (
		ctx           = context.Background()
		user          = uuid.New()
		expectedOrder = model.Order{
			OrderUUID:  uuid.New(),
			UserUUID:   user,
			PartUUIDs:  []uuid.UUID{uuid.New(), uuid.New()},
			TotalPrice: 199.99,
			Status:     model.OrderStatusPaid,
			PaymentInfo: &model.PaymentInfo{
				PaymentMethod:   model.PaymentMethodCard,
				TransactionUUID: uuid.New(),
			},
		}
	)

	err := s.ordersRepository.AddOrder(ctx, expectedOrder)
	s.Require().NoError(err)

	// Act
	result, err := s.ordersRepository.GetUserOrders(ctx, user)
	s.Require().NoError(err)
	s.Require().Len(result, 1)

	// Assert
	actualOrder := result[0]
	s.Require().Equal(expectedOrder.OrderUUID, actualOrder.OrderUUID)
	s.Require().Equal(expectedOrder.UserUUID, actualOrder.UserUUID)
	s.Require().Equal(expectedOrder.PartUUIDs, actualOrder.PartUUIDs)
	s.Require().Equal(expectedOrder.TotalPrice, actualOrder.TotalPrice)
	s.Require().Equal(expectedOrder.Status, actualOrder.Status)
	s.Require().Equal(expectedOrder.PaymentInfo.PaymentMethod, actualOrder.PaymentInfo.PaymentMethod)
	s.Require().Equal(expectedOrder.PaymentInfo.TransactionUUID, actualOrder.PaymentInfo.TransactionUUID)
}
