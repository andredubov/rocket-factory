package tests

import (
	"context"
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
)

// TestAddOrder_Success verifies successful order creation through the service layer.
// Tests that a valid order with all required fields can be properly processed
// and stored in the repository without errors.
func (s *OrdersServiceSuite) TestAddOrder_Success() {
	// Arrange
	ctx := context.Background()
	order := model.Order{
		OrderUUID: uuid.New(),
		UserUUID:  uuid.New(),
		PartUUIDs: []uuid.UUID{uuid.New(), uuid.New()},
		Status:    model.OrderStatusPending,
	}

	repoOrder := converter.OrderToRepoModel(order)
	s.ordersRepository.On("AddOrder", ctx, repoOrder).Return(nil)

	// Act
	err := s.ordersService.AddOrder(ctx, order)

	// Assert
	s.NoError(err)
}

// TestAddOrder_InvalidStatus verifies validation of order status during creation.
// Tests that orders with invalid status values are rejected with the appropriate
// ErrInvalidOrderStatus error before reaching the repository.
func (s *OrdersServiceSuite) TestAddOrder_InvalidStatus() {
	// Arrange
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			UserUUID:  uuid.New(),
			PartUUIDs: []uuid.UUID{uuid.New()},
			Status:    "INVALID_STATUS",
		}
		expectedError = model.ErrInvalidOrderStatus
	)

	repoOrder := converter.OrderToRepoModel(order)
	s.ordersRepository.On("AddOrder", ctx, repoOrder).Return(expectedError)

	// Act
	err := s.ordersService.AddOrder(ctx, order)

	// Assert
	s.Error(err)
	s.Equal(err, model.ErrInvalidOrderStatus)
}

// TestAddOrder_InvalidPaymentMethod verifies payment method validation.
// Tests that orders with invalid payment methods are rejected with
// ErrInvalidPaymentMethod error, ensuring only valid payment options are accepted.
func (s *OrdersServiceSuite) TestAddOrder_InvalidPaymentMethod() {
	// Arrange
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			UserUUID:  uuid.New(),
			PartUUIDs: []uuid.UUID{uuid.New()},
			Status:    model.OrderStatusPending,
			PaymentInfo: &model.PaymentInfo{
				PaymentMethod: "INVALID_METHOD",
			},
		}
		expectedError = model.ErrInvalidPaymentMethod
	)

	repoOrder := converter.OrderToRepoModel(order)
	s.ordersRepository.On("AddOrder", ctx, repoOrder).Return(expectedError)

	// Act
	err := s.ordersService.AddOrder(ctx, order)

	// Assert
	s.Error(err)
	s.Equal(err, model.ErrInvalidPaymentMethod)
}

// TestAddOrder_RepositoryError verifies proper error propagation from repository.
// Tests that repository-level errors (like duplicate orders) are correctly
// propagated through the service layer with their original error type.
func (s *OrdersServiceSuite) TestAddOrder_RepositoryError() {
	// Arrange
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			UserUUID:  uuid.New(),
			PartUUIDs: []uuid.UUID{uuid.New()},
			Status:    model.OrderStatusPending,
		}
		expectedError = model.ErrOrderAlreadyExists
	)

	repoOrder := converter.OrderToRepoModel(order)
	s.ordersRepository.On("AddOrder", ctx, repoOrder).Return(expectedError)

	// Act
	err := s.ordersService.AddOrder(ctx, order)

	// Assert
	s.Error(err)
	s.Equal(err, expectedError)
}

// TestAddOrder_WithGofakeit performs randomized testing of order creation.
// Generates multiple test cases with random valid data to verify the service
// handles various valid input combinations correctly. Includes random payment
// information generation to test optional field handling.
func (s *OrdersServiceSuite) TestAddOrder_WithGofakeit() {
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

			repoOrder := converter.OrderToRepoModel(order)
			s.ordersRepository.On("AddOrder", ctx, repoOrder).Return(nil)

			// Act
			err := s.ordersService.AddOrder(ctx, order)

			// Assert
			s.NoError(err)
		})
	}
}

// TestAddOrder_EmptyPartUUIDs verifies validation of orders with empty part lists.
// Tests that orders without any parts are rejected with an appropriate error,
// enforcing the business rule that orders must contain at least one part.
func (s *OrdersServiceSuite) TestAddOrder_EmptyPartUUIDs() {
	// Arrange
	var (
		ctx   = context.Background()
		order = model.Order{
			OrderUUID: uuid.New(),
			UserUUID:  uuid.New(),
			PartUUIDs: []uuid.UUID{},
			Status:    model.OrderStatusPending,
		}
		expectedError = errors.New("at least one part required")
	)

	repoOrder := converter.OrderToRepoModel(order)
	s.ordersRepository.On("AddOrder", ctx, repoOrder).Return(expectedError)

	// Act
	err := s.ordersService.AddOrder(ctx, order)

	// Assert
	s.Require().Error(err)
	s.Require().Equal(err, expectedError)
}
