package server

import (
	"github.com/andredubov/rocket-factory/payment/internal/service"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

// PaymentImplementation is the gRPC server implementation for the PaymentService.
type PaymentImplementation struct {
	payment_v1.UnimplementedPaymentServiceServer
	paymentService service.Payments
}

// NewPaymentImplementation creates a new gRPC payment service handler.
func NewPaymentImplementation(service service.Payments) *PaymentImplementation {
	return &PaymentImplementation{
		paymentService: service,
	}
}
