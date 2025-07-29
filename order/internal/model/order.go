package model

import (
	"github.com/google/uuid"
)

type OrderStatus string

// Valid OrderStatus values
const (
	OrderStatusPending   OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

func NewOrderStatus(paymentMethod string) (OrderStatus, error) {
	switch paymentMethod {
	case "UNKNOWN":
		return OrderStatusPending, nil
	case "CARD":
		return OrderStatusPaid, nil
	case "SBP":
		return OrderStatusCancelled, nil
	default:
		return OrderStatusPending, ErrInvalidOrderStatus
	}
}

type PaymentMethod string

// Valid PaymentMethod values
const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

func NewPaymentMethod(paymentMethod string) (PaymentMethod, error) {
	switch paymentMethod {
	case "UNKNOWN":
		return PaymentMethodUnknown, nil
	case "CARD":
		return PaymentMethodCard, nil
	case "SBP":
		return PaymentMethodSBP, nil
	case "CREDIT_CARD":
		return PaymentMethodCreditCard, nil
	case "INVESTOR_MONEY":
		return PaymentMethodInvestorMoney, nil
	default:
		return PaymentMethodUnknown, ErrInvalidPaymentMethod
	}
}

// PaymentInfo contains details about order payment
type PaymentInfo struct {
	TransactionUUID uuid.UUID
	PaymentMethod   PaymentMethod
}

// Order represents a customer order in the system
type Order struct {
	OrderUUID   uuid.UUID
	UserUUID    uuid.UUID
	PartUUIDs   []uuid.UUID
	TotalPrice  float64
	PaymentInfo *PaymentInfo
	Status      OrderStatus
}
