package orders

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/order/internal/converter"
	"github.com/andredubov/rocket-factory/order/internal/metrics"
	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	"github.com/andredubov/rocket-factory/platform/pkg/tracing"
)

// AddOrder creates a new order
func (s *ordersService) CreateOrder(ctx context.Context, order *model.Order) error {
	ctx, span := tracing.StartSpan(ctx, "order.create",
		trace.WithAttributes(
			attribute.String("user.uuid", order.UserUUID.String()),
			attribute.Int("parts.count", len(order.PartUUIDs)),
		),
	)
	defer span.End()

	// Валидация
	if len(order.PartUUIDs) == 0 {
		logger.Error(ctx, "❌ [OrderService] no parts in the order", zap.Error(model.ErrOrderHasNoParts))
		span.RecordError(model.ErrOrderHasNoParts)
		span.SetStatus(codes.Error, model.ErrOrderHasNoParts.Error())
		return model.ErrOrderHasNoParts
	}

	partFilter := converter.OrderToPartFilter(*order)

	parts, err := s.inventoryClient.ListParts(ctx, partFilter)
	if err != nil {
		logger.Error(ctx, "❌ [OrderService] failed to fetch the order parts", zap.Error(model.ErrInvalidPartFilter))
		span.RecordError(model.ErrInvalidPartFilter)
		span.SetStatus(codes.Error, model.ErrInvalidPartFilter.Error())
		return model.ErrInvalidPartFilter
	}

	order.Status = model.OrderStatusPending
	order.OrderUUID = uuid.New()
	order.TotalPrice = calculateTotalPrice(parts)

	metrics.IncOrdersTotal(ctx)
	metrics.AddRevenue(ctx, order.TotalPrice, "RUB")

	err = s.ordersRepository.AddOrder(ctx, *order)
	if err != nil {
		logger.Error(ctx, "❌ [OrderService] failed to create the order into database", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetAttributes(
		attribute.String("order.uuid", order.OrderUUID.String()),
	)
	span.SetStatus(codes.Ok, "order created successfully")

	return nil
}

// calculateTotalPrice calculates and return order total price
func calculateTotalPrice(parts []model.Part) float64 {
	total := decimal.NewFromFloat(0)
	for _, part := range parts {
		total = total.Add(decimal.NewFromFloat(part.Price))
	}
	result, _ := total.Round(2).Float64()

	return result
}
