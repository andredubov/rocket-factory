package tests

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// TestGetOrder_Success verifies successful retrieval of an order through the service layer.
// Tests that a valid order UUID returns the corresponding order with all fields properly
// converted from the repository model to the domain model.
func (s *OrdersServiceSuite) TestGetOrder_Success() {
	// Arrange
	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
		repoOrder = &repoModel.Order{
			OrderUUID: orderUUID,
			UserUUID:  uuid.New(),
			Status:    repoModel.OrderStatus(model.OrderStatusPaid),
		}
	)

	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(repoOrder, nil)

	// Act
	order, err := s.ordersService.GetOrder(ctx, orderUUID)

	// Assert
	s.NoError(err)
	s.Equal(converter.OrderToModel(*repoOrder), order)
}

// TestGetOrder_NotFound verifies proper handling of requests for non-existent orders.
// Tests that the service returns ErrOrderNotFound when attempting to retrieve
// an order that doesn't exist, maintaining proper error handling semantics.
func (s *OrdersServiceSuite) TestGetOrder_NotFound() {
	// Arrange
	var (
		ctx         = context.Background()
		orderUUID   = uuid.New()
		expectedErr = model.ErrOrderNotFound
	)

	s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(nil, expectedErr)

	// Act
	result, err := s.ordersService.GetOrder(ctx, orderUUID)

	// Assert
	s.Nil(result)
	s.Error(err)
	s.Equal(err, expectedErr)
}

// TestGetOrder_EmptyUUID verifies validation of empty UUID inputs.
// Tests that attempts to retrieve orders with empty UUIDs are properly rejected
// with ErrOrderNotFound, ensuring robust input validation.
func (s *OrdersServiceSuite) TestGetOrder_EmptyUUID() {
	// Arrange
	var (
		ctx         = context.Background()
		emptyUUID   = uuid.Nil
		expectedErr = model.ErrOrderNotFound
	)

	s.ordersRepository.On("GetOrder", ctx, emptyUUID).Return(nil, expectedErr)

	// Act
	result, err := s.ordersService.GetOrder(ctx, emptyUUID)

	// Assert
	s.Nil(result)
	s.Error(err)
	s.Equal(err, expectedErr)
}

// TestGetOrder_WithGofakeit performs randomized testing of order retrieval.
// Generates multiple test cases with varied order data and success/failure conditions
// to verify the service handles diverse query scenarios correctly.
func (s *OrdersServiceSuite) TestGetOrder_WithGofakeit() {
	// Generate 5 random test cases
	for i := 0; i < 5; i++ {
		s.Run(gofakeit.BeerName(), func() {
			// Arrange
			var (
				ctx       = context.Background()
				orderUUID = uuid.New()
				// Generate random order data
				repoOrder = &repoModel.Order{
					OrderUUID: orderUUID,
					UserUUID:  uuid.New(),
					PartUUIDs: []uuid.UUID{uuid.New(), uuid.New()},
					Status: repoModel.OrderStatus(gofakeit.RandomString([]string{
						string(repoModel.OrderStatusPending),
						string(repoModel.OrderStatusPaid),
						string(repoModel.OrderStatusCancelled),
					})),
					TotalPrice: gofakeit.Float64Range(10, 1000),
				}
			)

			// Randomly add payment info (50% chance)
			if gofakeit.Bool() {
				repoOrder.PaymentInfo = &repoModel.PaymentInfo{
					PaymentMethod: repoModel.PaymentMethod(gofakeit.RandomString([]string{
						string(repoModel.PaymentMethodCard),
						string(repoModel.PaymentMethodSBP),
						string(repoModel.PaymentMethodCreditCard),
						string(repoModel.PaymentMethodInvestorMoney),
					})),
					TransactionUUID: uuid.New(),
				}
			}

			expectedOrder := converter.OrderToModel(*repoOrder)

			// Randomly decide if test should succeed or fail
			if gofakeit.Bool() {
				// Success case
				s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(repoOrder, nil)

				// Act
				result, err := s.ordersService.GetOrder(ctx, orderUUID)

				// Assert
				s.NoError(err)
				s.Equal(expectedOrder, result)
			} else {
				// Failure case
				expectedErr := model.ErrOrderNotFound
				s.ordersRepository.On("GetOrder", ctx, orderUUID).Return(nil, expectedErr)

				// Act
				result, err := s.ordersService.GetOrder(ctx, orderUUID)

				// Assert
				s.Error(err)
				s.Nil(result)
				s.Equal(expectedErr, err)
			}
		})
	}
}
