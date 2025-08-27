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

// CreateOrder обрабатывает запрос на создание нового заказа.
func (i *OrderImplementation) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	order := converter.OrderFromCreateOrderRequest(req)

	err := i.ordersService.CreateOrder(ctx, &order)
	if err != nil {
		if errors.Is(err, model.ErrOrderHasNoParts) {
			return &order_v1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "at least one part required",
			}, nil
		}

		if errors.Is(err, model.ErrInvalidPartFilter) {
			return &order_v1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("invalid part %v", err),
			}, nil
		}

		if errors.Is(err, model.ErrOrderAlreadyExists) {
			return &order_v1.ConflictError{
				Code:    http.StatusConflict,
				Message: "order already exists",
			}, nil
		}

		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return converter.OrderToCreateOrderResponse(order), nil
}
