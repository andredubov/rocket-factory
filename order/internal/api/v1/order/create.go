package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/andredubov/rocket-factory/order/internal/converter"
	"github.com/andredubov/rocket-factory/order/internal/model"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// CreateOrder обрабатывает запрос на создание нового заказа.
func (i *OrderImplementation) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	order := converter.OrderFromCreateOrderRequest(req)
	order.Status = model.OrderStatusPending
	order.OrderUUID = uuid.New()

	// Валидация
	if len(order.PartUUIDs) == 0 {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "at least one part required",
		}, nil
	}

	// Проверка наличия деталей и расчет стоимости
	total := decimal.NewFromFloat(0)
	for _, partUuid := range order.PartUUIDs {
		inventoryRequest := inventory_v1.GetPartRequest{Uuid: partUuid.String()}
		inventoryResponse, err := i.inventoryClient.GetPart(ctx, &inventoryRequest)
		if err != nil {
			return &order_v1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("invalid part %s: %v", partUuid, err),
			}, nil
		}

		total = total.Add(decimal.NewFromFloat(inventoryResponse.GetPart().GetPrice()))
	}

	order.TotalPrice, _ = total.Round(2).Float64()

	// Сохранение
	if err := i.ordersService.AddOrder(ctx, order); err != nil {
		if errors.Is(err, model.ErrOrderAlreadyExists) {
			return &order_v1.ConflictError{
				Code:    http.StatusConflict,
				Message: "order already exists",
			}, nil
		}
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return converter.OrderToCreateOrderResponse(order), nil
}
