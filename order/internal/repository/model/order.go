package model

import "github.com/gofrs/uuid"

// OrderStatus represents the possible states of an order
type OrderStatus string

// IsValid checks if the OrderStatus has a valid value
// Returns:
//   - true if status is one of the defined constants
//   - false otherwise
func (os OrderStatus) IsValid() bool {
	switch os {
	case OrderStatusPending, OrderStatusPaid, OrderStatusCancelled:
		return true
	default:
		return false
	}
}

// Valid OrderStatus values
const (
	// Order is awaiting payment
	OrderStatusPending OrderStatus = "PENDING_PAYMENT"
	// Order has been paid
	OrderStatusPaid OrderStatus = "PAID"
	// Order was cancelled
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// PaymentMethod represents the payment options available for orders
type PaymentMethod string

// IsValid checks if the PaymentMethod has a valid value
// Returns:
//   - true if method is one of the defined constants
//   - false otherwise
func (pm PaymentMethod) IsValid() bool {
	switch pm {
	case PaymentMethodCard, PaymentMethodSBP, PaymentMethodCreditCard, PaymentMethodBankTransfer:
		return true
	default:
		return false
	}
}

// Valid PaymentMethod values
const (
	// Standard debit/credit card payment
	PaymentMethodCard PaymentMethod = "CARD"
	// Faster Payments System (Russian payment system)
	PaymentMethodSBP PaymentMethod = "SBP"
	// Credit card payment
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	// Bank transfer payment
	PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
)

// PaymentInfo contains details about order payment
type PaymentInfo struct {
	// Unique identifier for the payment transaction
	TransactionUUID uuid.UUID
	// Method used for payment
	PaymentMethod PaymentMethod
}

// Order represents a customer order in the system
type Order struct {
	// Unique identifier for the order
	OrderUUID uuid.UUID
	// Identifier of the user who placed the order
	UserUUID uuid.UUID
	// List of part identifiers included in the order
	PartUUIDs []uuid.UUID
	// Total price of the order in minor units (e.g., cents)
	TotalPrice int64
	// Payment details (nullable)
	PaymentInfo *PaymentInfo
	// Current status of the order
	Status OrderStatus
}
