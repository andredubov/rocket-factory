package orders

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/andredubov/rocket-factory/order/internal/converter"
	"github.com/andredubov/rocket-factory/order/internal/metrics"
	"github.com/andredubov/rocket-factory/order/internal/model"
)

// AddOrder creates a new order
func (s *ordersService) CreateOrder(ctx context.Context, order *model.Order) error {
	// Валидация
	if len(order.PartUUIDs) == 0 {
		return model.ErrOrderHasNoParts
	}

	partFilter := converter.OrderToPartFilter(*order)

	parts, err := s.inventoryClient.ListParts(ctx, partFilter)
	if err != nil {
		return model.ErrInvalidPartFilter
	}

	order.Status = model.OrderStatusPending
	order.OrderUUID = uuid.New()
	order.TotalPrice = calculateTotalPrice(parts)

	metrics.IncOrdersTotal(ctx)
	metrics.AddRevenue(ctx, order.TotalPrice, "RUB")

	return s.ordersRepository.AddOrder(ctx, *order)
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
