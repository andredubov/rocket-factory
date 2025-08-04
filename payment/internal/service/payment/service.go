package payment

import api "github.com/andredubov/rocket-factory/payment/internal/api/v1/payment"

// paymentService is the internal implementation of the Payments interface.
// It handles core payment processing logic and acts as the main entry point
// for payment operations in the domain layer.
type paymentService struct{}

// NewService creates and returns a new instance of the payment service.
func NewService() api.PaymentService {
	return &paymentService{}
}
