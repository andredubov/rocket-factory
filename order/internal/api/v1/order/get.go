package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/andredubov/rocket-factory/order/internal/converter"
	"github.com/andredubov/rocket-factory/order/internal/model"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

// GetOrderByUuid возвращает информацию о заказе по его UUID.
func (i *OrderImplementation) GetOrderByUuid(ctx context.Context, params order_v1.GetOrderByUuidParams) (order_v1.GetOrderByUuidRes, error) {
	// Проверяем не отменен ли контекст перед началом работы
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context cancelled before processing: %w", err)
	}

	// Получаем заказ из репозитория
	order, err := i.ordersService.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			// Возвращаем структурированную ошибку API для случая "не найдено"
			return &order_v1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("order with UUID %s not found", params.OrderUUID),
			}, nil
		}
		// Все остальные ошибки возвращаем как внутренние ошибки сервера
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return converter.OrderToGetOrderResponse(order)
}
