package payment

import (
	"github.com/andredubov/rocket-factory/order/internal/client/grpc"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

type paymentClient struct {
	generatedClient payment_v1.PaymentServiceClient
}

// NewClient...
func NewClient(client payment_v1.PaymentServiceClient) grpc.PaymentClient {
	return &paymentClient{
		generatedClient: client,
	}
}
