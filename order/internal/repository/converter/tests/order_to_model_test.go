package tests

import (
	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// TestOrderToModel_FullConversion verifies complete conversion of a repository order to domain model.
// Tests that all fields including nested payment information are correctly mapped,
// validating the complete transformation of complex order structures.
func (s *OrdersRepoConverterSuite) TestOrderToModel_FullConversion() {
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
	s.Require().Equal(repoOrder.OrderUUID, result.OrderUUID)
	s.Require().Equal(repoOrder.UserUUID, result.UserUUID)
	s.Require().Equal(repoOrder.PartUUIDs, result.PartUUIDs)
	s.Require().Equal(repoOrder.TotalPrice, result.TotalPrice)
	s.Require().Equal(model.OrderStatus(repoOrder.Status), result.Status)

	s.Require().NotNil(result.PaymentInfo)
	s.Require().Equal(repoOrder.PaymentInfo.TransactionUUID, result.PaymentInfo.TransactionUUID)
	s.Require().Equal(model.PaymentMethod(repoOrder.PaymentInfo.PaymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToModel_NoPaymentInfo verifies proper handling of orders without payment information.
// Tests that nil payment info in the repository model correctly converts to nil
// in the domain model, ensuring optional field handling.
func (s *OrdersRepoConverterSuite) TestOrderToModel_NoPaymentInfo() {
	// Arrange
	repoOrder := repoModel.Order{
		OrderUUID: uuid.New(),
		Status:    repoModel.OrderStatusPending,
		// PaymentInfo is nil
	}

	// Act
	result := converter.OrderToModel(repoOrder)

	// Assert
	s.Require().Equal(repoOrder.OrderUUID, result.OrderUUID)
	s.Require().Equal(model.OrderStatus(repoOrder.Status), result.Status)
	s.Require().Nil(result.PaymentInfo)
}

// TestOrderToModel_EmptyPartUUIDs verifies correct conversion of orders with empty part lists.
// Tests that empty part UUID slices are properly preserved during conversion,
// maintaining data structure integrity.
func (s *OrdersRepoConverterSuite) TestOrderToModel_EmptyPartUUIDs() {
	// Arrange
	repoOrder := repoModel.Order{
		OrderUUID: uuid.New(),
		PartUUIDs: []uuid.UUID{}, // Empty slice
		Status:    repoModel.OrderStatusCancelled,
	}

	// Act
	result := converter.OrderToModel(repoOrder)

	// Assert
	s.Empty(result.PartUUIDs)
	s.Equal(model.OrderStatusCancelled, result.Status)
}

// testRepoOrderStatusConversion is a helper function for testing status enum conversions.
// Provides consistent testing of repository-to-domain status value mapping
// for different order status values.
func (s *OrdersRepoConverterSuite) testRepoOrderStatusConversion(status repoModel.OrderStatus) {
	repoOrder := repoModel.Order{
		OrderUUID: uuid.New(),
		Status:    status,
	}

	result := converter.OrderToModel(repoOrder)
	s.Equal(model.OrderStatus(status), result.Status)
}

// TestOrderToModel_StatusPending verifies Pending status conversion.
// Tests specific mapping of repository's Pending status to domain model equivalent.
func (s *OrdersRepoConverterSuite) TestOrderToModel_StatusPending() {
	s.testRepoOrderStatusConversion(repoModel.OrderStatusPending)
}

// TestOrderToModel_StatusPaid verifies Paid status conversion.
// Tests specific mapping of repository's Paid status to domain model equivalent.
func (s *OrdersRepoConverterSuite) TestOrderToModel_StatusPaid() {
	s.testRepoOrderStatusConversion(repoModel.OrderStatusPaid)
}

// TestOrderToModel_StatusCancelled verifies Cancelled status conversion.
// Tests specific mapping of repository's Cancelled status to domain model equivalent.
func (s *OrdersRepoConverterSuite) TestOrderToModel_StatusCancelled() {
	s.testRepoOrderStatusConversion(repoModel.OrderStatusCancelled)
}

// testRepoOrderPaymentMethodConversion is a helper for payment method enum tests.
// Provides consistent testing of payment method value conversions between
// repository and domain models.
func (s *OrdersRepoConverterSuite) testRepoOrderPaymentMethodConversion(paymentMethod repoModel.PaymentMethod) {
	repoOrder := repoModel.Order{
		OrderUUID: uuid.New(),
		PaymentInfo: &repoModel.PaymentInfo{
			PaymentMethod: paymentMethod,
		},
	}

	result := converter.OrderToModel(repoOrder)
	s.NotNil(result.PaymentInfo)
	s.Equal(model.PaymentMethod(paymentMethod), result.PaymentInfo.PaymentMethod)
}

// TestOrderToModel_PaymentMethodUnknown verifies Unknown payment method conversion.
// Tests specific mapping of repository's Unknown payment method to domain model.
func (s *OrdersRepoConverterSuite) TestOrderToModel_PaymentMethodUnknown() {
	s.testRepoOrderPaymentMethodConversion(repoModel.PaymentMethodUnknown)
}

// TestOrderToModel_PaymentMethodCard verifies Card payment method conversion.
// Tests specific mapping of repository's Card payment method to domain model.
func (s *OrdersRepoConverterSuite) TestOrderToModel_PaymentMethodCard() {
	s.testRepoOrderPaymentMethodConversion(repoModel.PaymentMethodCard)
}

// TestOrderToModel_PaymentMethodSBP verifies SBP payment method conversion.
// Tests specific mapping of repository's SBP payment method to domain model.
func (s *OrdersRepoConverterSuite) TestOrderToModel_PaymentMethodSBP() {
	s.testRepoOrderPaymentMethodConversion(repoModel.PaymentMethodSBP)
}

// TestOrderToModel_PaymentMethodCreditCard verifies Credit Card payment method conversion.
// Tests specific mapping of repository's Credit Card payment method to domain model.
func (s *OrdersRepoConverterSuite) TestOrderToModel_PaymentMethodCreditCard() {
	s.testRepoOrderPaymentMethodConversion(repoModel.PaymentMethodCreditCard)
}

// TestOrderToModel_PaymentMethodInvestorMoney verifies Investor Money payment method conversion.
// Tests specific mapping of repository's Investor Money payment method to domain model.
func (s *OrdersRepoConverterSuite) TestOrderToModel_PaymentMethodInvestorMoney() {
	s.testRepoOrderPaymentMethodConversion(repoModel.PaymentMethodInvestorMoney)
}
