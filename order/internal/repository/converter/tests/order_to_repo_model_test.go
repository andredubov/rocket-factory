package tests

import (
	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// TestOrderToRepoModel_FullConversion verifies complete conversion of a domain order to repository model.
// Tests that all fields including nested payment information are correctly mapped,
// validating the complete transformation of complex order structures.
func (s *OrdersRepoConverterSuite) TestOrderToRepoModel_FullConversion() {
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
	s.Require().Equal(order.OrderUUID, result.OrderUUID)
	s.Require().Equal(order.UserUUID, result.UserUUID)
	s.Require().Equal(order.PartUUIDs, result.PartUUIDs)
	s.Require().Equal(order.TotalPrice, result.TotalPrice)
	s.Require().Equal(repoModel.OrderStatus(order.Status), result.Status)

	s.Require().NotNil(result.PaymentInfo)
	s.Require().Equal(order.PaymentInfo.TransactionUUID, result.PaymentInfo.TransactionUUID)
	s.Require().Equal(repoModel.PaymentMethod(order.PaymentInfo.PaymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToRepoModel_NoPaymentInfo verifies proper handling of orders without payment information.
// Tests that nil payment info in the domain model correctly converts to nil
// in the repository model, ensuring proper handling of optional fields.
func (s *OrdersRepoConverterSuite) TestOrderToRepoModel_NoPaymentInfo() {
	// Arrange
	order := model.Order{
		OrderUUID: uuid.New(),
		Status:    model.OrderStatusPending,
		// PaymentInfo is nil
	}

	// Act
	result := converter.OrderToRepoModel(order)

	// Assert
	s.Require().Equal(order.OrderUUID, result.OrderUUID)
	s.Require().Equal(repoModel.OrderStatus(order.Status), result.Status)
	s.Require().Nil(result.PaymentInfo)
}

// TestOrderToRepoModel_EmptyPartUUIDs verifies correct conversion of orders with empty part lists.
// Tests that empty part UUID slices are properly preserved during conversion,
// maintaining data structure integrity between layers.
func (s *OrdersRepoConverterSuite) TestOrderToRepoModel_EmptyPartUUIDs() {
	// Arrange
	order := model.Order{
		OrderUUID: uuid.New(),
		PartUUIDs: []uuid.UUID{}, // Empty slice
		Status:    model.OrderStatusCancelled,
	}

	// Act
	result := converter.OrderToRepoModel(order)

	// Assert
	s.Require().Empty(result.PartUUIDs)
	s.Require().Equal(repoModel.OrderStatusCancelled, result.Status)
}

// TestOrderToRepoModel_AllStatuses verifies comprehensive status enum conversion.
// Tests all possible order status values to ensure complete and correct mapping
// between domain and repository status representations.
func (s *OrdersRepoConverterSuite) TestOrderToRepoModel_AllStatuses() {
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
		s.Require().Equal(repoModel.OrderStatus(tc.status), result.Status)
	}
}

// testStatusConversion is a helper function for testing status enum conversions.
// Provides consistent testing of domain-to-repository status value mapping
// for different order status values.
func (s *OrdersRepoConverterSuite) testStatusConversion(status model.OrderStatus) {
	order := model.Order{
		OrderUUID: uuid.New(),
		Status:    status,
	}
	result := converter.OrderToRepoModel(order)
	s.Require().Equal(repoModel.OrderStatus(status), result.Status)
}

// TestOrderToRepoModel_StatusPending verifies Pending status conversion.
// Tests specific mapping of domain's Pending status to repository model equivalent.
func (s *OrdersRepoConverterSuite) TestOrderToRepoModel_StatusPending() {
	s.testStatusConversion(model.OrderStatusPending)
}

// TestOrderToRepoModel_StatusPaid verifies Paid status conversion.
// Tests specific mapping of domain's Paid status to repository model equivalent.
func (s *OrdersRepoConverterSuite) TestOrderToRepoModel_StatusPaid() {
	s.testStatusConversion(model.OrderStatusPaid)
}

// TestOrderToRepoModel_StatusCancelled verifies Cancelled status conversion.
// Tests specific mapping of domain's Cancelled status to repository model equivalent.
func (s *OrdersRepoConverterSuite) TestOrderToRepoModel_StatusCancelled() {
	s.testStatusConversion(model.OrderStatusCancelled)
}

// testPaymentMethodConversion is a helper for payment method enum tests.
// Provides consistent testing of payment method value conversions between
// domain and repository models.
func (s *OrdersRepoConverterSuite) testPaymentMethodConversion(paymentMethod model.PaymentMethod) {
	order := model.Order{
		OrderUUID: uuid.New(),
		PaymentInfo: &model.PaymentInfo{
			PaymentMethod: paymentMethod,
		},
	}

	result := converter.OrderToRepoModel(order)
	s.Require().NotNil(result.PaymentInfo)
	s.Require().Equal(repoModel.PaymentMethod(paymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToRepoModel_PaymentMethodUnknown verifies Unknown payment method conversion.
// Tests specific mapping of domain's Unknown payment method to repository model.
func (s *OrdersRepoConverterSuite) TestOrderToRepoModel_PaymentMethodUnknown() {
	s.testPaymentMethodConversion(model.PaymentMethodUnknown)
}

// TestOrderToRepoModel_PaymentMethodCard verifies Card payment method conversion.
// Tests specific mapping of domain's Card payment method to repository model
func (s *OrdersRepoConverterSuite) TestOrderToRepoModel_PaymentMethodCard() {
	s.testPaymentMethodConversion(model.PaymentMethodCard)
}

// TestOrderToRepoModel_PaymentMethodSBP verifies SBP payment method conversion.
// Tests specific mapping of domain's SBP payment method to repository model.
func (s *OrdersRepoConverterSuite) TestOrderToRepoModel_PaymentMethodSBP() {
	s.testPaymentMethodConversion(model.PaymentMethodSBP)
}

// TestOrderToRepoModel_PaymentMethodCreditCard verifies Credit Card payment method conversion.
// Tests specific mapping of domain's Credit Card payment method to repository model.
func (s *OrdersRepoConverterSuite) TestOrderToRepoModel_PaymentMethodCreditCard() {
	s.testPaymentMethodConversion(model.PaymentMethodCreditCard)
}

// TestOrderToRepoModel_PaymentMethodInvestorMoney verifies Investor Money payment method conversion.
// Tests specific mapping of domain's Investor Money payment method to repository model.
func (s *OrdersRepoConverterSuite) TestOrderToRepoModel_PaymentMethodInvestorMoney() {
	s.testPaymentMethodConversion(model.PaymentMethodInvestorMoney)
}
