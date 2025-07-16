package converter

import (
	"fmt"

	"github.com/andredubov/rocket-factory/order/internal/model"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

// OrderFromCreateOrderRequest converts a CreateOrderRequest protobuf message to a domain Order model.
func OrderFromCreateOrderRequest(req *order_v1.CreateOrderRequest) model.Order {
	return model.Order{
		UserUUID:  req.GetUserUUID(),
		PartUUIDs: req.GetPartUuids(),
	}
}

// OrderToCreateOrderResponse converts a domain Order model to a CreateOrderResponse protobuf message.
func OrderToCreateOrderResponse(order model.Order) *order_v1.CreateOrderResponse {
	return &order_v1.CreateOrderResponse{
		OrderUUID:  order_v1.NewOptUUID(order.OrderUUID),
		TotalPrice: order_v1.NewOptFloat64(order.TotalPrice),
	}
}

// OrderToPayOrderResponse converts a domain Order model to a PayOrderResponse protobuf message.
func OrderToPayOrderResponse(order *model.Order) *order_v1.PayOrderResponse {
	return &order_v1.PayOrderResponse{
		TransactionUUID: order_v1.NewOptUUID(order.PaymentInfo.TransactionUUID),
	}
}

// OrderToGetOrderResponse converts a domain Order model to a PayOrderResponse protobuf message.
func OrderToGetOrderResponse(order *model.Order) (*order_v1.GetOrderResponse, error) {
	response := &order_v1.GetOrderResponse{
		OrderUUID:  order.OrderUUID,                    // UUID заказа
		UserUUID:   order.UserUUID,                     // UUID пользователя
		PartUuids:  order.PartUUIDs,                    // Список UUID деталей в заказе
		TotalPrice: order.TotalPrice,                   // Общая стоимость заказа
		Status:     order_v1.OrderStatus(order.Status), // Текущий статус заказа
	}

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
