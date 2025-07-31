package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/andredubov/rocket-factory/order/internal/model"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

// CancelOrder обрабатывает запрос на отмену заказа.
func (i *OrderImplementation) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	// отменяем заказ
	err := i.ordersService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "order not found",
			}, nil
		}

		if errors.Is(err, model.ErrOrderAlreadyPaid) {
			// Если заказ уже оплачен - возвращаем ошибку 409
			return &order_v1.ConflictError{
				Code:    http.StatusConflict,
				Message: "order has been paid and cannot be cancelled",
			}, nil
		}

		if errors.Is(err, model.ErrOrderAlreadyCancelled) {
			// Если заказ уже отменен - возвращаем ошибку 409
			return &order_v1.ConflictError{
				Code:    http.StatusConflict,
				Message: "order is already cancelled",
			}, nil
		}

		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return &order_v1.CancelOrderNoContent{}, nil
}
