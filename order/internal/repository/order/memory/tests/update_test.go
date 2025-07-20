package tests

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// TestUpdateOrder_Success verifies successful order updates in the repository.
// Tests the complete flow: adding an order, updating its status and payment info,
// and verifying the changes persist. Validates field-level updates includin status transitions and payment information.
func (s *OrdersRepositorySuite) TestUpdateOrder_Success() {
	// Setup
	var (
		ctx           = context.Background()
		originalOrder = model.Order{
			OrderUUID:  uuid.New(),
			UserUUID:   uuid.New(),
			Status:     model.OrderStatusPending,
			TotalPrice: 100.0,
		}
	)

	// Add initial order
	err := s.ordersRepository.AddOrder(ctx, originalOrder)
	s.Require().NoError(err)

	// Prepare update
	updatedOrder := originalOrder
	updatedOrder.Status = model.OrderStatusPaid
	updatedOrder.TotalPrice = 150.0
	updatedOrder.PaymentInfo = &model.PaymentInfo{
		PaymentMethod:   model.PaymentMethodCard,
		TransactionUUID: uuid.New(),
	}

	// Test
	err = s.ordersRepository.UpdateOrder(ctx, updatedOrder)

	// Verify
	s.NoError(err)

	// Verify changes were applied
	retrieved, err := s.ordersRepository.GetOrder(ctx, originalOrder.OrderUUID)
	s.Require().NoError(err)
	s.Require().Equal(updatedOrder.Status, retrieved.Status)
	s.Require().Equal(updatedOrder.TotalPrice, retrieved.TotalPrice)
	s.Require().Equal(updatedOrder.PaymentInfo.PaymentMethod, retrieved.PaymentInfo.PaymentMethod)
}

// TestUpdateOrder_NotFound verifies proper handling of update attempts for non-existent orders.
// Tests that the repository returns ErrOrderNotFound when attempting to update
// an order that doesn't exist, ensuring data integrity.
func (s *OrdersRepositorySuite) TestUpdateOrder_NotFound() {
	// Setup
	ctx := context.Background()
	nonExistentOrder := model.Order{
		OrderUUID: uuid.New(),
		Status:    model.OrderStatusPending,
	}

	// Test
	err := s.ordersRepository.UpdateOrder(ctx, nonExistentOrder)

	// Verify
	s.Require().Equal(err, repository.ErrOrderNotFoundWith(nonExistentOrder.OrderUUID))
}

// TestUpdateOrder_InvalidStatus verifies validation of order status during updates.
// Tests that invalid status values are rejected with ErrInvalidOrderStatus,
// maintaining business rule enforcement.
func (s *OrdersRepositorySuite) TestUpdateOrder_InvalidStatus() {
	// Setup
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			Status:    "INVALID_STATUS",
		}
	)

	// Test
	err := s.ordersRepository.UpdateOrder(ctx, order)

	// Verify
	s.Require().Equal(err, repository.ErrInvalidOrderStatusWith("INVALID_STATUS"))
}

// TestUpdateOrder_InvalidPaymentMethod verifies payment method validation during updates.
// Tests that orders with invalid payment methods are rejected with ErrInvalidPaymentMethod,
// ensuring payment processing requirements are met.
func (s *OrdersRepositorySuite) TestUpdateOrder_InvalidPaymentMethod() {
	// Setup
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			Status:    model.OrderStatusPaid,
			PaymentInfo: &model.PaymentInfo{
				PaymentMethod: "INVALID_METHOD",
			},
		}
	)

	// Test
	err := s.ordersRepository.UpdateOrder(ctx, order)

	// Verify
	s.Require().Equal(err, repository.ErrInvalidPaymentMethodWith("INVALID_METHOD"))
}

// TestUpdateOrder_ConcurrentAccess verifies thread-safe update behavior under race conditions.
// Tests that concurrent updates to the same order are handled correctly, with
// the repository applying one of the updates while maintaining data consistency.
func (s *OrdersRepositorySuite) TestUpdateOrder_ConcurrentAccess() {
	// Setup
	var (
		ctx           = context.Background()
		originalOrder = model.Order{
			OrderUUID:   uuid.New(),
			Status:      model.OrderStatusPending,
			PaymentInfo: &model.PaymentInfo{PaymentMethod: model.PaymentMethodCard},
		}
	)

	// Add initial order
	err := s.ordersRepository.AddOrder(ctx, originalOrder)
	s.Require().NoError(err)

	// Prepare updates
	update1 := originalOrder
	update1.TotalPrice = 300.0

	update2 := originalOrder
	update2.TotalPrice = 200.0

	// Test concurrent updates
	var wg sync.WaitGroup
	wg.Add(2)

	var err1, err2 error
	go func() {
		defer wg.Done()
		err1 = s.ordersRepository.UpdateOrder(ctx, update1)
	}()
	go func() {
		defer wg.Done()
		err2 = s.ordersRepository.UpdateOrder(ctx, update2)
	}()

	wg.Wait()

	// Verify both succeeded
	s.Require().True((err1 == nil && err2 == nil))

	// Verify order was updated
	retrieved, err := s.ordersRepository.GetOrder(ctx, originalOrder.OrderUUID)
	s.NoError(err)

	// Check which update was applied
	if err1 == nil {
		s.Require().Equal(update1.Status, retrieved.Status)
	} else {
		s.Require().Equal(update2.TotalPrice, retrieved.TotalPrice)
	}
}
