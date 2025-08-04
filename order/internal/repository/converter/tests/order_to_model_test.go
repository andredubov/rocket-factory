package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// TestOrderToModel_FullConversion verifies complete conversion of a repository order to domain model.
func TestOrderToModel_FullConversion(t *testing.T) {
	// Arrange
	repoOrder := repoModel.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice: 199.99,
		Status:     repoModel.OrderStatusPaid,
		PaymentInfo: &repoModel.PaymentInfo{
			TransactionUUID: uuid.New(),
			PaymentMethod:   repoModel.PaymentMethodCard,
		},
	}

	// Act
	result := converter.OrderToModel(repoOrder)

	// Assert
	require.Equal(t, repoOrder.OrderUUID, result.OrderUUID)
	require.Equal(t, repoOrder.UserUUID, result.UserUUID)
	require.Equal(t, repoOrder.PartUUIDs, result.PartUUIDs)
	require.Equal(t, repoOrder.TotalPrice, result.TotalPrice)
	require.Equal(t, model.OrderStatus(repoOrder.Status), result.Status)

	require.NotNil(t, result.PaymentInfo)
	require.Equal(t, repoOrder.PaymentInfo.TransactionUUID, result.PaymentInfo.TransactionUUID)
	require.Equal(t, model.PaymentMethod(repoOrder.PaymentInfo.PaymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToModel_NoPaymentInfo verifies proper handling of orders without payment information.
func TestOrderToModel_NoPaymentInfo(t *testing.T) {
	// Arrange
	repoOrder := repoModel.Order{
		OrderUUID: uuid.New(),
		Status:    repoModel.OrderStatusPending,
		// PaymentInfo is nil
	}

	// Act
	result := converter.OrderToModel(repoOrder)

	// Assert
	require.Equal(t, repoOrder.OrderUUID, result.OrderUUID)
	require.Equal(t, model.OrderStatus(repoOrder.Status), result.Status)
	require.Nil(t, result.PaymentInfo)
}

// TestOrderToModel_EmptyPartUUIDs verifies correct conversion of orders with empty part lists.
func TestOrderToModel_EmptyPartUUIDs(t *testing.T) {
	// Setup
	repoOrder := repoModel.Order{
		OrderUUID: uuid.New(),
		PartUUIDs: []uuid.UUID{},
		Status:    repoModel.OrderStatusCancelled,
	}

	// Act
	result := converter.OrderToModel(repoOrder)

	// Assert
	require.Empty(t, result.PartUUIDs)
	require.Equal(t, model.OrderStatusCancelled, result.Status)
}

// testRepoOrderStatusConversion is a helper function for testing status enum conversions.
func testRepoOrderStatusConversion(status repoModel.OrderStatus) *model.Order {
	repoOrder := repoModel.Order{
		OrderUUID: uuid.New(),
		Status:    status,
	}
	return converter.OrderToModel(repoOrder)
}

// TestOrderToModel_StatusPending verifies Pending status conversion.
func TestOrderToModel_StatusPending(t *testing.T) {
	status := repoModel.OrderStatusPending
	result := testRepoOrderStatusConversion(status)
	require.Equal(t, model.OrderStatus(status), result.Status)
}

// TestOrderToModel_StatusPaid verifies Paid status conversion.
func TestOrderToModel_StatusPaid(t *testing.T) {
	status := repoModel.OrderStatusPaid
	result := testRepoOrderStatusConversion(status)
	require.Equal(t, model.OrderStatus(status), result.Status)
}

// TestOrderToModel_StatusCancelled verifies Cancelled status conversion.
func TestOrderToModel_StatusCancelled(t *testing.T) {
	status := repoModel.OrderStatusCancelled
	result := testRepoOrderStatusConversion(status)
	require.Equal(t, model.OrderStatus(status), result.Status)
}

// testRepoOrderPaymentMethodConversion is a helper for payment method enum tests.
func testRepoOrderPaymentMethodConversion(paymentMethod repoModel.PaymentMethod) *model.Order {
	repoOrder := repoModel.Order{
		OrderUUID: uuid.New(),
		PaymentInfo: &repoModel.PaymentInfo{
			PaymentMethod: paymentMethod,
		},
	}
	return converter.OrderToModel(repoOrder)
}

// TestOrderToModel_PaymentMethodUnknown verifies Unknown payment method conversion.
func TestOrderToModel_PaymentMethodUnknown(t *testing.T) {
	paymentMethod := repoModel.PaymentMethodUnknown
	result := testRepoOrderPaymentMethodConversion(paymentMethod)
	require.NotNil(t, result)
	require.Equal(t, model.PaymentMethod(paymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToModel_PaymentMethodCard verifies Card payment method conversion.
func TestOrderToModel_PaymentMethodCard(t *testing.T) {
	paymentMethod := repoModel.PaymentMethodCard
	result := testRepoOrderPaymentMethodConversion(paymentMethod)
	require.NotNil(t, result)
	require.Equal(t, model.PaymentMethod(paymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToModel_PaymentMethodSBP verifies SBP payment method conversion.
func TestOrderToModel_PaymentMethodSBP(t *testing.T) {
	paymentMethod := repoModel.PaymentMethodSBP
	result := testRepoOrderPaymentMethodConversion(paymentMethod)
	require.NotNil(t, result)
	require.Equal(t, model.PaymentMethod(paymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToModel_PaymentMethodCreditCard verifies Credit Card payment method conversion.
func TestOrderToModel_PaymentMethodCreditCard(t *testing.T) {
	paymentMethod := repoModel.PaymentMethodCreditCard
	result := testRepoOrderPaymentMethodConversion(paymentMethod)
	require.NotNil(t, result)
	require.Equal(t, model.PaymentMethod(paymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToModel_PaymentMethodInvestorMoney verifies Investor Money payment method conversion.
func TestOrderToModel_PaymentMethodInvestorMoney(t *testing.T) {
	paymentMethod := repoModel.PaymentMethodInvestorMoney
	result := testRepoOrderPaymentMethodConversion(paymentMethod)
	require.NotNil(t, result)
	require.Equal(t, model.PaymentMethod(paymentMethod), result.PaymentInfo.PaymentMethod)
}
