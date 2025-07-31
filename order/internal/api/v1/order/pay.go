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

// PayOrder обрабатывает запрос на оплату заказа.
func (i *OrderImplementation) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
	paymentMethod := converter.OrderPaymentMethodFromRequest(req)
	// Оплачиваем заказ
	order, err := i.ordersService.PayOrder(ctx, params.OrderUUID, paymentMethod)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "order not found",
			}, nil
		}

		if errors.Is(err, model.ErrInvalidOrderStatus) {
			return &order_v1.ConflictError{
				Code:    http.StatusConflict,
				Message: "order is not in pending status",
			}, nil
		}

		if errors.Is(err, model.ErrInvalidPaymentMethod) {
			return &order_v1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "invalid payment method",
			}, nil
		}

		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return converter.OrderToPayOrderResponse(order), nil
}
