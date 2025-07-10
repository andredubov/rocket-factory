package model

import "github.com/google/uuid"

type OrderStatus string

// Valid OrderStatus values
const (
	OrderStatusPending   OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// IsValid checks if the OrderStatus has a valid value
func (os OrderStatus) IsValid() bool {
	switch os {
	case OrderStatusPending, OrderStatusPaid, OrderStatusCancelled:
		return true
	default:
		return false
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

// IsValid checks if the PaymentMethod has a valid value
func (pm PaymentMethod) IsValid() bool {
	switch pm {
	case PaymentMethodUnknown, PaymentMethodCard, PaymentMethodSBP, PaymentMethodCreditCard, PaymentMethodInvestorMoney:
		return true
	default:
		return false
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
