package payment

import (
	"github.com/andredubov/rocket-factory/order/internal/service"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

// paymentClient is a client for interacting with the Payment Service.
// It wraps the auto-generated gRPC client to provide payment-specific functionality.
type paymentClient struct {
	generatedClient payment_v1.PaymentServiceClient
}

// NewClient creates a new instance of paymentClient.
func NewClient(client payment_v1.PaymentServiceClient) service.PaymentClient {
	return &paymentClient{
		generatedClient: client,
	}
}
