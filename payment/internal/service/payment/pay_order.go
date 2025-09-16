package payment

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/payment/internal/model"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	"github.com/andredubov/rocket-factory/platform/pkg/tracing"
)

// Create implements the Payments interface by generating a new payment transaction.
func (p *paymentService) Create(ctx context.Context, payment model.Payment) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "payment.process",
		trace.WithAttributes(attribute.String("order.id", payment.OrderUuid)),
	)
	defer span.End()

	transactionUUID := uuid.New().String()
	logger.Info(ctx, "✅ [PaymentService] Payment was successful", zap.String("transaction_uuid:", transactionUUID))

	span.SetAttributes(attribute.String("payment.transaction.uuid", transactionUUID))
	span.SetStatus(codes.Ok, "payment succeeded")

	return transactionUUID, nil
}
