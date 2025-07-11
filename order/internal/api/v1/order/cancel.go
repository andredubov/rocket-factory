package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/repository/model"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

// CancelOrder обрабатывает запрос на отмену заказа.
func (i *OrderImplementation) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	// Получаем заказ из репозитория
	order, err := i.ordersRepository.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "order not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Проверяем статус заказа
	switch order.Status {
	case model.OrderStatusPaid:
		// Если заказ уже оплачен - возвращаем ошибку 409
		return &order_v1.ConflictError{
			Code:    http.StatusConflict,
			Message: "order has been paid and cannot be cancelled",
		}, nil
	case model.OrderStatusCancelled:
		// Если заказ уже отменен - возвращаем ошибку 409
		return &order_v1.ConflictError{
			Code:    http.StatusConflict,
			Message: "order is already cancelled",
		}, nil
	}

	// Меняем статус на Cancelled для заказов в статусе Pending
	order.Status = model.OrderStatusCancelled

	// Обновляем заказ в репозитории
	if err := i.ordersRepository.UpdateOrder(ctx, *order); err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "order not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return &order_v1.CancelOrderNoContent{}, nil
}
