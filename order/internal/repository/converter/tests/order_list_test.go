package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// TestOrdersToModel_EmptySlice verifies correct handling of empty order slices during conversion.
// Tests that converting an empty repository order slice returns an empty domain model slice,
// ensuring proper handling of zero-length inputs.
func TestOrdersToModel_EmptySlice(t *testing.T) {
	// Arrange
	var repoOrders []repoModel.Order

	// Act
	result := converter.OrdersToModel(repoOrders)

	// Assert
	require.Empty(t, result)
}

// TestOrdersToModel_SingleOrder verifies accurate conversion of a single order.
// Tests that all fields of a single repository order are correctly mapped to the
// corresponding domain model fields, including UUID and status conversion.
func TestOrdersToModel_SingleOrder(t *testing.T) {
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
	require.Len(t, result, 1)
	require.Equal(t, repoOrders[0].OrderUUID, result[0].OrderUUID)
	require.Equal(t, model.OrderStatus(repoOrders[0].Status), result[0].Status)
}

// TestOrdersToModel_MultipleOrders verifies batch conversion of multiple orders with different states.
// Tests comprehensive conversion of an order slice containing various statuses and
// payment information, validating field-level mapping for each order in the collection.
func TestOrdersToModel_MultipleOrders(t *testing.T) {
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
	require.Len(t, result, 3)

	// Verify first order
	require.Equal(t, repoOrders[0].OrderUUID, result[0].OrderUUID)
	require.Equal(t, model.OrderStatusPending, result[0].Status)
	require.Nil(t, result[0].PaymentInfo)

	// Verify second order
	require.Equal(t, repoOrders[1].OrderUUID, result[1].OrderUUID)
	require.Equal(t, model.OrderStatusPaid, result[1].Status)
	require.NotNil(t, result[1].PaymentInfo)
	require.Equal(t, model.PaymentMethodCard, result[1].PaymentInfo.PaymentMethod)

	// Verify third order
	require.Equal(t, repoOrders[2].OrderUUID, result[2].OrderUUID)
	require.Equal(t, model.OrderStatusCancelled, result[2].Status)
}
