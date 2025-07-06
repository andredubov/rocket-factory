package model

import payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"

// Payment represents a payment transaction in the domain layer.
// It contains essential information linking a user's payment method to a specific order.
type Payment struct {
	// UserID identifies the user making the payment.
	// Note: Currently sourced from OrderUuid in the proto request,
	UserID string

	// OrderID references the associated order for this payment.
	// This should match the order system's unique identifier.
	OrderID string

	// PaymentMethod specifies how the payment is being processed.
	// This is converted directly from the gRPC enum PaymentMethod.
	PaymentMethod payment_v1.PaymentMethod
}
