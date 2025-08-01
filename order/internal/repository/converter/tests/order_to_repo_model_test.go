package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// TestOrderToRepoModel_FullConversion verifies complete conversion of a domain order to repository model.
func TestOrderToRepoModel_FullConversion(t *testing.T) {
	// Arrange
	order := model.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice: 199.99,
		Status:     model.OrderStatusPaid,
		PaymentInfo: &model.PaymentInfo{
			TransactionUUID: uuid.New(),
			PaymentMethod:   model.PaymentMethodCard,
		},
	}

	// Act
	result := converter.OrderToRepoModel(order)

	// Assert
	require.Equal(t, order.OrderUUID, result.OrderUUID)
	require.Equal(t, order.UserUUID, result.UserUUID)
	require.Equal(t, order.PartUUIDs, result.PartUUIDs)
	require.Equal(t, order.TotalPrice, result.TotalPrice)
	require.Equal(t, repoModel.OrderStatus(order.Status), result.Status)

	require.NotNil(t, result.PaymentInfo)
	require.Equal(t, order.PaymentInfo.TransactionUUID, result.PaymentInfo.TransactionUUID)
	require.Equal(t, repoModel.PaymentMethod(order.PaymentInfo.PaymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToRepoModel_NoPaymentInfo verifies proper handling of orders without payment information.
func TestOrderToRepoModel_NoPaymentInfo(t *testing.T) {
	// Arrange
	order := model.Order{
		OrderUUID: uuid.New(),
		Status:    model.OrderStatusPending,
		// PaymentInfo is nil
	}

	// Act
	result := converter.OrderToRepoModel(order)

	// Assert
	require.Equal(t, order.OrderUUID, result.OrderUUID)
	require.Equal(t, repoModel.OrderStatus(order.Status), result.Status)
	require.Nil(t, result.PaymentInfo)
}

// TestOrderToRepoModel_EmptyPartUUIDs verifies correct conversion of orders with empty part lists.
func TestOrderToRepoModel_EmptyPartUUIDs(t *testing.T) {
	// Arrange
	order := model.Order{
		OrderUUID: uuid.New(),
		PartUUIDs: []uuid.UUID{}, // Empty slice
		Status:    model.OrderStatusCancelled,
	}

	// Act
	result := converter.OrderToRepoModel(order)

	// Assert
	require.Empty(t, result.PartUUIDs)
	require.Equal(t, repoModel.OrderStatusCancelled, result.Status)
}

// TestOrderToRepoModel_AllStatuses verifies comprehensive status enum conversion.
func TestOrderToRepoModel_AllStatuses(t *testing.T) {
	// Test all possible order status values
	testCases := []struct {
		name   string
		status model.OrderStatus
	}{
		{"Pending", model.OrderStatusPending},
		{"Paid", model.OrderStatusPaid},
		{"Cancelled", model.OrderStatusCancelled},
	}

	for _, tc := range testCases {
		// Arrange
		order := model.Order{
			OrderUUID: uuid.New(),
			Status:    tc.status,
		}
		// Act
		result := converter.OrderToRepoModel(order)
		// Assert
		require.Equal(t, repoModel.OrderStatus(tc.status), result.Status)
	}
}

// testStatusConversion is a helper function for testing status enum conversions.
// Provides consistent testing of domain-to-repository status value mapping
// for different order status values.
func testStatusConversion(status model.OrderStatus) repoModel.Order {
	order := model.Order{
		OrderUUID: uuid.New(),
		Status:    status,
	}
	return converter.OrderToRepoModel(order)
}

// TestOrderToRepoModel_StatusPending verifies Pending status conversion.
// Tests specific mapping of domain's Pending status to repository model equivalent.
func TestOrderToRepoModel_StatusPending(t *testing.T) {
	status := model.OrderStatusPending
	result := testStatusConversion(status)
	require.Equal(t, repoModel.OrderStatus(status), result.Status)
}

// TestOrderToRepoModel_StatusPaid verifies Paid status conversion.
// Tests specific mapping of domain's Paid status to repository model equivalent.
func TestOrderToRepoModel_StatusPaid(t *testing.T) {
	status := model.OrderStatusPaid
	result := testStatusConversion(status)
	require.Equal(t, repoModel.OrderStatus(status), result.Status)
}

// TestOrderToRepoModel_StatusCancelled verifies Cancelled status conversion.
// Tests specific mapping of domain's Cancelled status to repository model equivalent.
func TestOrderToRepoModel_StatusCancelled(t *testing.T) {
	status := model.OrderStatusCancelled
	result := testStatusConversion(status)
	require.Equal(t, repoModel.OrderStatus(status), result.Status)
}

// testPaymentMethodConversion is a helper for payment method enum tests.
// Provides consistent testing of payment method value conversions between
// domain and repository models.
func testPaymentMethodConversion(paymentMethod model.PaymentMethod) repoModel.Order {
	order := model.Order{
		OrderUUID: uuid.New(),
		PaymentInfo: &model.PaymentInfo{
			PaymentMethod: paymentMethod,
		},
	}
	return converter.OrderToRepoModel(order)
}

// TestOrderToRepoModel_PaymentMethodUnknown verifies Unknown payment method conversion.
// Tests specific mapping of domain's Unknown payment method to repository model.
func TestOrderToRepoModel_PaymentMethodUnknown(t *testing.T) {
	paymentMethod := model.PaymentMethodUnknown
	result := testPaymentMethodConversion(paymentMethod)
	require.NotNil(t, result.PaymentInfo)
	require.Equal(t, repoModel.PaymentMethod(paymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToRepoModel_PaymentMethodCard verifies Card payment method conversion.
// Tests specific mapping of domain's Card payment method to repository model
func TestOrderToRepoModel_PaymentMethodCard(t *testing.T) {
	paymentMethod := model.PaymentMethodCard
	result := testPaymentMethodConversion(paymentMethod)
	require.NotNil(t, result.PaymentInfo)
	require.Equal(t, repoModel.PaymentMethod(paymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToRepoModel_PaymentMethodSBP verifies SBP payment method conversion.
// Tests specific mapping of domain's SBP payment method to repository model.
func TestOrderToRepoModel_PaymentMethodSBP(t *testing.T) {
	paymentMethod := model.PaymentMethodSBP
	result := testPaymentMethodConversion(paymentMethod)
	require.NotNil(t, result.PaymentInfo)
	require.Equal(t, repoModel.PaymentMethod(paymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToRepoModel_PaymentMethodCreditCard verifies Credit Card payment method conversion.
// Tests specific mapping of domain's Credit Card payment method to repository model.
func TestOrderToRepoModel_PaymentMethodCreditCard(t *testing.T) {
	paymentMethod := model.PaymentMethodCreditCard
	result := testPaymentMethodConversion(paymentMethod)
	require.NotNil(t, result.PaymentInfo)
	require.Equal(t, repoModel.PaymentMethod(paymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToRepoModel_PaymentMethodInvestorMoney verifies Investor Money payment method conversion.
// Tests specific mapping of domain's Investor Money payment method to repository model.
func TestOrderToRepoModel_PaymentMethodInvestorMoney(t *testing.T) {
	paymentMethod := model.PaymentMethodInvestorMoney
	result := testPaymentMethodConversion(paymentMethod)
	require.NotNil(t, result.PaymentInfo)
	require.Equal(t, repoModel.PaymentMethod(paymentMethod), result.PaymentInfo.PaymentMethod)
}
