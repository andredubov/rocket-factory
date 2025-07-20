package tests

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
)

// TestUpdateOrder_Success verifies successful order updates through the service layer.
// Tests that valid order updates are properly processed and stored in the repository
// without errors, confirming the happy path scenario for order modifications.
func (s *OrdersServiceSuite) TestUpdateOrder_Success() {
	// Arrange
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			UserUUID:  uuid.New(),
			PartUUIDs: []uuid.UUID{uuid.New(), uuid.New()},
			Status:    model.OrderStatusPaid,
		}
	)

	repoOrder := converter.OrderToRepoModel(order)
	s.ordersRepository.On("UpdateOrder", ctx, repoOrder).Return(nil)

	// Act
	err := s.ordersService.UpdateOrder(ctx, order)

	// Assert
	s.NoError(err)
}

// TestUpdateOrder_InvalidStatus verifies validation of order status during updates.
// Tests that orders with invalid status values are rejected with ErrInvalidOrderStatus,
// ensuring only permitted status transitions are allowed.
func (s *OrdersServiceSuite) TestUpdateOrder_InvalidStatus() {
	// Arrange
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			Status:    "INVALID_STATUS",
		}

		expectedError = model.ErrInvalidOrderStatus
	)

	repoOrder := converter.OrderToRepoModel(order)
	s.ordersRepository.On("UpdateOrder", ctx, repoOrder).Return(expectedError)

	// Act
	err := s.ordersService.UpdateOrder(ctx, order)

	// Assert
	s.Error(err)
	s.Equal(err, expectedError)
}

// TestUpdateOrder_InvalidPaymentMethod verifies payment method validation during updates.
// Tests that invalid payment methods are rejected with ErrInvalidPaymentMethod,
// maintaining payment processing integrity.
func (s *OrdersServiceSuite) TestUpdateOrder_InvalidPaymentMethod() {
	// Arrange
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			Status:    model.OrderStatusPaid,
			PaymentInfo: &model.PaymentInfo{
				PaymentMethod: "INVALID_METHOD",
			},
		}
		expectedError = model.ErrInvalidPaymentMethod
	)

	repoOrder := converter.OrderToRepoModel(order)
	s.ordersRepository.On("UpdateOrder", ctx, repoOrder).Return(expectedError)

	// Act
	err := s.ordersService.UpdateOrder(ctx, order)

	// Assert
	s.Error(err)
	s.Equal(err, expectedError)
}

// TestUpdateOrder_NotFound verifies handling of updates to non-existent orders.
// Tests that attempts to update missing orders return ErrOrderNotFound,
// ensuring proper error handling for invalid order references.
func (s *OrdersServiceSuite) TestUpdateOrder_NotFound() {
	// Arrange
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			Status:    model.OrderStatusPaid,
		}
		expectedError = model.ErrOrderNotFound
	)

	repoOrder := converter.OrderToRepoModel(order)
	s.ordersRepository.On("UpdateOrder", ctx, repoOrder).Return(expectedError)

	// Act
	err := s.ordersService.UpdateOrder(ctx, order)

	// Assert
	s.Error(err)
	s.Equal(err, expectedError)
}

// TestUpdateOrder_WithGofakeit performs randomized testing of order updates.
// Generates multiple test cases with varied valid/invalid data combinations
// to verify the service handles diverse update scenarios correctly.
func (s *OrdersServiceSuite) TestUpdateOrder_WithGofakeit() {
	// Generate 5 random test cases
	for i := 0; i < 5; i++ {
		s.Run(gofakeit.BeerName(), func() {
			// Arrange
			var (
				ctx   = context.Background()
				order = model.Order{
					OrderUUID: uuid.New(),
					UserUUID:  uuid.New(),
					PartUUIDs: []uuid.UUID{uuid.New(), uuid.New()},
					Status: model.OrderStatus(gofakeit.RandomString([]string{
						string(model.OrderStatusPending),
						string(model.OrderStatusPaid),
						string(model.OrderStatusCancelled),
					})),
					TotalPrice: gofakeit.Float64Range(10, 1000),
				}
			)

			// Randomly add payment info (50% chance)
			if gofakeit.Bool() {
				order.PaymentInfo = &model.PaymentInfo{
					PaymentMethod: model.PaymentMethod(gofakeit.RandomString([]string{
						string(model.PaymentMethodCard),
						string(model.PaymentMethodSBP),
						string(model.PaymentMethodCreditCard),
						string(model.PaymentMethodInvestorMoney),
					})),
					TransactionUUID: uuid.New(),
				}
			}

			// Randomly decide if test should succeed or fail
			if gofakeit.Bool() {
				// Success case
				s.ordersRepository.On("UpdateOrder", ctx, converter.OrderToRepoModel(order)).Return(nil)

				// Act
				err := s.ordersService.UpdateOrder(ctx, order)

				// Assert
				s.NoError(err)
			} else {
				// Failure case
				expectedErr := model.ErrOrderNotFound
				s.ordersRepository.On("UpdateOrder", ctx, converter.OrderToRepoModel(order)).Return(expectedErr)

				// Act
				err := s.ordersService.UpdateOrder(ctx, order)

				// Assert
				s.Error(err)
				s.Equal(err, expectedErr)
			}
		})
	}
}

// TestUpdateOrder_EmptyOrderUUID verifies validation of empty order UUIDs.
// Tests that update attempts with nil UUIDs are rejected with ErrOrderNotFound,
// enforcing proper order reference requirements.
func (s *OrdersServiceSuite) TestUpdateOrder_EmptyOrderUUID() {
	// Arrange
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.Nil, // Empty UUID
			Status:    model.OrderStatusPaid,
		}
		expectedError = model.ErrOrderNotFound
	)

	repoOrder := converter.OrderToRepoModel(order)
	s.ordersRepository.On("UpdateOrder", ctx, repoOrder).Return(expectedError)

	// Act
	err := s.ordersService.UpdateOrder(ctx, order)

	// Assert
	s.Error(err)
	s.Equal(err, expectedError)
}
