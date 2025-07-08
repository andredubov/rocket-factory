package order

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/repository/model"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

func (i *Implementation) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
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

	// Создание запроса в платежный сервис
	paymentRequest := &payment_v1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: ConvertModelPaymentMethodToProto(order.PaymentInfo.PaymentMethod),
	}

	// Вызов платежного сервиса
	paymentResponse, err := i.paymentClient.PayOrder(ctx, paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("payment service error: %w", err)
	}

	// Парсинг UUID транзакции
	transactionUUID, err := uuid.Parse(paymentResponse.GetTransactionUuid())
	if err != nil {
		return nil, fmt.Errorf("invalid transaction uuid: %w", err)
	}

	// Обновление информации о заказе
	order.PaymentInfo.TransactionUUID = transactionUUID
	order.Status = model.OrderStatusPaid

	// Сохранение обновленного заказа
	if err := i.ordersRepository.UpdateOrder(ctx, *order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Формирование ответа
	response := &order_v1.PayOrderResponse{
		TransactionUUID: order_v1.NewOptUUID(transactionUUID),
	}

	return response, nil
}

// ConvertModelPaymentMethodToProto конвертирует model.PaymentMethod в payment_v1.PaymentMethod
func ConvertModelPaymentMethodToProto(method model.PaymentMethod) payment_v1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
