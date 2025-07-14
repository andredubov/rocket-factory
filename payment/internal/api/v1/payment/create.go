package server

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/andredubov/rocket-factory/payment/internal/converter"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

// PayOrder handles gRPC request for order payment.
func (i *PaymentImplementation) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	// Convert gRPC request to domain model
	payment := converter.PaymentFromRequest(req)

	// Create payment through domain service
	uuid, err := i.paymentService.Create(ctx, payment)
	if err != nil {
		log.Printf("payment creation failed: %v", err)
		return nil, status.Errorf(codes.Internal, "payment creation failed: %v", err)
	}

	// Log successful transaction
	log.Printf("Оплата прошла успешно, transaction_uuid: %s\n", uuid)

	// Return response with transaction ID
	return &payment_v1.PayOrderResponse{
		TransactionUuid: uuid,
	}, nil
}
