package tests

import (
	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// TestOrdersToModel_EmptySlice verifies correct handling of empty order slices during conversion.
// Tests that converting an empty repository order slice returns an empty domain model slice,
// ensuring proper handling of zero-length inputs.
func (s *OrdersRepoConverterSuite) TestOrdersToModel_EmptySlice() {
	// Arrange
	var repoOrders []repoModel.Order

	// Act
	result := converter.OrdersToModel(repoOrders)

	// Assert
	s.Empty(result)
}

// TestOrdersToModel_SingleOrder verifies accurate conversion of a single order.
// Tests that all fields of a single repository order are correctly mapped to the
// corresponding domain model fields, including UUID and status conversion.
func (s *OrdersRepoConverterSuite) TestOrdersToModel_SingleOrder() {
	// Arrange
	repoOrders := []repoModel.Order{
		{
			OrderUUID: uuid.New(),
			UserUUID:  uuid.New(),
			Status:    repoModel.OrderStatusPaid,
		},
	}

	// Act
	result := converter.OrdersToModel(repoOrders)

	// Assert
	s.Len(result, 1)
	s.Equal(repoOrders[0].OrderUUID, result[0].OrderUUID)
	s.Equal(model.OrderStatus(repoOrders[0].Status), result[0].Status)
}

// TestOrdersToModel_MultipleOrders verifies batch conversion of multiple orders with different states.
// Tests comprehensive conversion of an order slice containing various statuses and
// payment information, validating field-level mapping for each order in the collection.
func (s *OrdersRepoConverterSuite) TestOrdersToModel_MultipleOrders() {
	// Arrange
	repoOrders := []repoModel.Order{
		{
			OrderUUID: uuid.New(),
			Status:    repoModel.OrderStatusPending,
		},
		{
			OrderUUID: uuid.New(),
			Status:    repoModel.OrderStatusPaid,
			PaymentInfo: &repoModel.PaymentInfo{
				PaymentMethod: repoModel.PaymentMethodCard,
			},
		},
		{
			OrderUUID: uuid.New(),
			Status:    repoModel.OrderStatusCancelled,
		},
	}

	// Act
	result := converter.OrdersToModel(repoOrders)

	// Assert
	s.Len(result, 3)

	// Verify first order
	s.Equal(repoOrders[0].OrderUUID, result[0].OrderUUID)
	s.Equal(model.OrderStatusPending, result[0].Status)
	s.Nil(result[0].PaymentInfo)

	// Verify second order
	s.Equal(repoOrders[1].OrderUUID, result[1].OrderUUID)
	s.Equal(model.OrderStatusPaid, result[1].Status)
	s.NotNil(result[1].PaymentInfo)
	s.Equal(model.PaymentMethodCard, result[1].PaymentInfo.PaymentMethod)

	// Verify third order
	s.Equal(repoOrders[2].OrderUUID, result[2].OrderUUID)
	s.Equal(model.OrderStatusCancelled, result[2].Status)
}
