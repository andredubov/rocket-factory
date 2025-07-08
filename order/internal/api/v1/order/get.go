package order

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/repository/model"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

func (i *Implementation) GetOrderByUuid(ctx context.Context, params order_v1.GetOrderByUuidParams) (order_v1.GetOrderByUuidRes, error) {
	// Проверяем не отменен ли контекст перед началом работы
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context cancelled before processing: %w", err)
	}

	// Получаем заказ из репозитория
	// В случае ошибки проверяем, является ли она ошибкой "не найдено"
	order, err := i.ordersRepository.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
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
	res := &order_v1.GetOrderResponse{
		OrderUUID:  order.OrderUUID,                    // UUID заказа
		UserUUID:   order.UserUUID,                     // UUID пользователя
		PartUuids:  order.PartUUIDs,                    // Список UUID деталей в заказе
		TotalPrice: order.TotalPrice,                   // Общая стоимость заказа
		Status:     order_v1.OrderStatus(order.Status), // Текущий статус заказа
	}

	// Если есть информация о платеже, добавляем ее в ответ
	if order.PaymentInfo != nil {
		// Добавляем UUID транзакции (может быть nil)
		res.TransactionUUID = order_v1.NewOptNilUUID(order.PaymentInfo.TransactionUUID)

		// Конвертируем метод оплаты из внутреннего формата в API-формат
		paymentMethod, err := convertToPaymentMethod(order.PaymentInfo.PaymentMethod)
		if err != nil {
			return nil, fmt.Errorf("payment method conversion error: %w", err)
		}
		res.PaymentMethod = order_v1.NewOptPaymentMethod(paymentMethod)
	}

	return res, nil
}

// convertToPaymentMethod конвертирует внутреннее представление метода оплаты
// в формат, используемый в API.
// Параметры:
//   - method: метод оплаты во внутреннем формате
//
// Возвращает:
//   - order_v1.PaymentMethod: метод оплаты в API-формате
//   - error: ошибка, если передан неизвестный метод оплаты
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
