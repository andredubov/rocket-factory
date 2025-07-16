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

	// Создаем базовый ответ с обязательными полями заказа
	response := converter.OrderToGetOrderResponse(order)

	// Если есть информация о платеже, добавляем ее в ответ
	if order.PaymentInfo != nil {
		response.TransactionUUID = order_v1.NewOptNilUUID(order.PaymentInfo.TransactionUUID)

		paymentMethod, err := convertToPaymentMethod(order.PaymentInfo.PaymentMethod)
		if err != nil {
			return nil, fmt.Errorf("payment method conversion error: %w", err)
		}
		response.PaymentMethod = order_v1.NewOptPaymentMethod(paymentMethod)
	}

	return response, nil
}

// convertToPaymentMethod конвертирует внутреннее представление метода оплаты в формат API.
func convertToPaymentMethod(method model.PaymentMethod) (order_v1.PaymentMethod, error) {
	switch method {
	case model.PaymentMethodCard:
		return order_v1.PaymentMethodCARD, nil
	case model.PaymentMethodSBP:
		return order_v1.PaymentMethodSBP, nil
	case model.PaymentMethodCreditCard:
		return order_v1.PaymentMethodCREDITCARD, nil
	case model.PaymentMethodInvestorMoney:
		return order_v1.PaymentMethodINVESTORMONEY, nil
	default:
		return "", fmt.Errorf("unsupported payment method: %s", method)
	}
}
