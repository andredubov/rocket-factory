package tests

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// TestGetUserOrders_Success verifies successful retrieval of a user's orders through the service layer.
// Tests that all orders for a specific user are properly returned and converted
// from repository to domain models, including nested payment information when present.
func (s *OrdersServiceSuite) TestGetUserOrders_Success() {
	// Arrange
	var (
		ctx        = context.Background()
		userUUID   = uuid.New()
		repoOrders = []repoModel.Order{
			{
				OrderUUID: uuid.New(),
				UserUUID:  userUUID,
				Status:    repoModel.OrderStatusPaid,
			},
			{
				OrderUUID: uuid.New(),
				UserUUID:  userUUID,
				Status:    repoModel.OrderStatusPending,
			},
		}
	)

	expectedOrders := converter.OrdersToModel(repoOrders)
	s.ordersRepository.On("GetUserOrders", ctx, userUUID).Return(repoOrders, nil)

	// Act
	result, err := s.ordersService.GetUserOrders(ctx, userUUID)

	// Assert
	s.NoError(err)
	s.Equal(expectedOrders, result)
}

// TestGetUserOrders_EmptyResult verifies proper handling when a user has no orders.
// Tests that an empty slice (not nil) is returned without errors when querying
// a user with no order history, ensuring correct empty state behavior.
func (s *OrdersServiceSuite) TestGetUserOrders_EmptyResult() {
	// Arrange
	var (
		ctx      = context.Background()
		userUUID = uuid.New()
	)

	s.ordersRepository.On("GetUserOrders", ctx, userUUID).Return([]repoModel.Order{}, nil)

	// Act
	result, err := s.ordersService.GetUserOrders(ctx, userUUID)

	// Assert
	s.NoError(err)
	s.Empty(result)
}

// TestGetUserOrders_RepositoryError verifies error propagation from the repository layer.
// Tests that repository-level errors are properly propagated through the service
// while maintaining the original error type and context.
func (s *OrdersServiceSuite) TestGetUserOrders_RepositoryError() {
	// Arrange
	var (
		ctx         = context.Background()
		userUUID    = uuid.New()
		expectedErr = model.ErrOrderNotFound
	)
	s.ordersRepository.On("GetUserOrders", ctx, userUUID).Return(nil, expectedErr)

	// Act
	result, err := s.ordersService.GetUserOrders(ctx, userUUID)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

// TestGetUserOrders_WithGofakeit performs randomized testing of order list retrieval.
// Generates multiple test cases with varied order data (including random payment info)
// and success/failure conditions to thoroughly test the service's query capabilities.
func (s *OrdersServiceSuite) TestGetUserOrders_WithGofakeit() {
	// Generate 5 random test cases
	for i := 0; i < 5; i++ {
		s.Run(gofakeit.BeerName(), func() {
			// Arrange
			var (
				ctx       = context.Background()
				userUUID  = uuid.New()
				numOrders = gofakeit.Number(1, 10)
			)

			// Generate random orders
			var repoOrders []repoModel.Order
			for j := 0; j < numOrders; j++ {
				order := repoModel.Order{
					OrderUUID: uuid.New(),
					UserUUID:  userUUID,
					Status: repoModel.OrderStatus(gofakeit.RandomString([]string{
						string(repoModel.OrderStatusPending),
						string(repoModel.OrderStatusPaid),
						string(repoModel.OrderStatusCancelled),
					})),
					TotalPrice: gofakeit.Float64Range(10, 1000),
				}

				// Randomly add payment info (50% chance)
				if gofakeit.Bool() {
					order.PaymentInfo = &repoModel.PaymentInfo{
						PaymentMethod: repoModel.PaymentMethod(gofakeit.RandomString([]string{
							string(repoModel.PaymentMethodCard),
							string(repoModel.PaymentMethodSBP),
							string(repoModel.PaymentMethodCreditCard),
							string(repoModel.PaymentMethodInvestorMoney),
						})),
						TransactionUUID: uuid.New(),
					}
				}
				repoOrders = append(repoOrders, order)
			}

			expectedOrders := converter.OrdersToModel(repoOrders)

			// Randomly decide if test should succeed or fail
			if gofakeit.Bool() {
				// Success case
				s.ordersRepository.On("GetUserOrders", ctx, userUUID).Return(repoOrders, nil)

				// Act
				result, err := s.ordersService.GetUserOrders(ctx, userUUID)

				// Assert
				s.NoError(err)
				s.Equal(expectedOrders, result)
			} else {
				// Failure case
				expectedErr := model.ErrOrderNotFound
				s.ordersRepository.On("GetUserOrders", ctx, userUUID).Return(nil, expectedErr)

				// Act
				result, err := s.ordersService.GetUserOrders(ctx, userUUID)

				// Assert
				s.Error(err)
				s.Nil(result)
				s.Equal(expectedErr, err)
			}
		})
	}
}

// TestGetUserOrders_EmptyUUID verifies input validation for empty user UUIDs.
// Tests that requests with nil UUIDs are properly rejected with ErrOrderNotFound,
// enforcing valid user reference requirements.
func (s *OrdersServiceSuite) TestGetUserOrders_EmptyUUID() {
	// Arrange
	var (
		ctx         = context.Background()
		emptyUUID   = uuid.Nil
		expectedErr = model.ErrOrderNotFound
	)

	s.ordersRepository.On("GetUserOrders", ctx, emptyUUID).Return(nil, expectedErr)

	// Act
	result, err := s.ordersService.GetUserOrders(ctx, emptyUUID)

	// Assert
	s.Nil(result)
	s.Error(err)
	s.Equal(err, expectedErr)
}
