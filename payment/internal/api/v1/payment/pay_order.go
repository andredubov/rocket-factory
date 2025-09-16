package api

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/andredubov/rocket-factory/payment/internal/converter"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

// PayOrder handles gRPC request for order payment.
func (i *PaymentImplementation) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	// Convert gRPC request to domain model
	payment := converter.PaymentFromRequest(req)

	// Create payment through domain service
	uuid, err := i.paymentService.Create(ctx, payment)
	if err != nil {
		logger.Error(ctx, "payment creation failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "payment creation failed: %v", err)
	}

	// Log successful transaction
	logger.Error(ctx, "Оплата прошла успешно", zap.String("transaction_uuid", uuid))

	// Return response with transaction ID
	return converter.TransactionUuidToResponse(uuid), nil
}
