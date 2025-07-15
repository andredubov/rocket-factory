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
	// Получаем заказ из репозитория
	order, err := i.ordersService.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "order not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Валидация статуса заказа
	if order.Status != model.OrderStatusPending {
		return &order_v1.ConflictError{
			Code:    http.StatusConflict,
			Message: "order is not in pending status",
		}, nil
	}

	// Валидация метода оплаты
	if !model.PaymentMethod(req.PaymentMethod).IsValid() {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "invalid payment method",
		}, nil
	}

	// Подготовка платежной информации
	order.PaymentInfo = &model.PaymentInfo{
		PaymentMethod: model.PaymentMethod(req.PaymentMethod),
	}

	// Вызов платежного сервиса
	transactionUUID, err := i.paymentClient.PayOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction uuid: %w", err)
	}

	// Обновление информации о заказе
	order.PaymentInfo.TransactionUUID = transactionUUID
	order.Status = model.OrderStatusPaid

	// Сохранение обновленного заказа
	if err := i.ordersService.UpdateOrder(ctx, *order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return converter.OrderToPayOrderResponse(order), nil
}
