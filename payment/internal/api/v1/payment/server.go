package api

import (
	"context"

	"github.com/andredubov/rocket-factory/payment/internal/model"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

// PaymentService defines the contract for payment processing operations.
type PaymentService interface {
	Create(ctx context.Context, payment model.Payment) (string, error)
}

// PaymentImplementation is the gRPC server implementation for the PaymentService.
type PaymentImplementation struct {
	payment_v1.UnimplementedPaymentServiceServer
	paymentService PaymentService
}

// NewPaymentImplementation creates a new gRPC payment service handler.
func NewPaymentImplementation(service PaymentService) *PaymentImplementation {
	return &PaymentImplementation{
		paymentService: service,
	}
}
