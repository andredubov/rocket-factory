package server

import (
	"github.com/andredubov/rocket-factory/payment/internal/service"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

// Implementation is the gRPC server implementation for the PaymentService.
// It wraps the domain-level payment service and adapts gRPC requests to domain calls.
//
// The struct embeds the generated UnimplementedPaymentServiceServer to ensure forward
// compatibility when adding new RPC methods.
type Implementation struct {
	payment_v1.UnimplementedPaymentServiceServer

	// paymentService is the underlying domain service that handles business logic.
	// All gRPC methods delegate to this service after request adaptation.
	paymentService service.Payments
}

// NewImplementation creates a new gRPC payment service handler.
//
// Parameters:
//
//	service - The domain payment service that contains business logic.
//	          This dependency injection allows for easier testing and
//	          separates transport concerns from business logic.
//
// Returns:
//
//	*Implementation - Configured gRPC service handler ready for registration.
func NewImplementation(service service.Payments) *Implementation {
	return &Implementation{
		paymentService: service,
	}
}
