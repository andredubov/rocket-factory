package metrics

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	serviceName = "order-service"

	ordersCreatedMetricName  = "orders_total"
	ordersCreatedDescription = "Общее количество созданных заказов"

	orderRevenueMetricName  = "orders_revenue_total"
	orderRevenueDescription = "Суммарная выручка от заказов"

	currencyUnit      = "currency"
	currencyAttribute = "currency"
)

var (
	ordersTotal        metric.Int64Counter   // ordersTotal - счетчик созданных заказов
	ordersRevenueTotal metric.Float64Counter // ordersRevenueTotal - суммарная выручка
)

// Init инициализирует все инструменты метрик с описаниями для production
func Init(_ context.Context) error {
	var err error

	meter := otel.Meter(serviceName)

	ordersTotal, err = meter.Int64Counter(
		ordersCreatedMetricName,
		metric.WithDescription(ordersCreatedDescription),
	)
	if err != nil {
		return err
	}

	ordersRevenueTotal, err = meter.Float64Counter(
		orderRevenueMetricName,
		metric.WithDescription(orderRevenueDescription),
		metric.WithUnit(currencyUnit),
	)
	if err != nil {
		return err
	}

	return nil
}

// IncOrdersTotal увеличивает счетчик созданных заказов
func IncOrdersTotal(ctx context.Context) {
	ordersTotal.Add(ctx, 1)
}

// AddRevenue добавляет выручку от заказа
func AddRevenue(ctx context.Context, amount float64, currency string) {
	ordersRevenueTotal.Add(ctx, amount,
		metric.WithAttributes(
			attribute.String(currencyAttribute, currency),
		),
	)
}
