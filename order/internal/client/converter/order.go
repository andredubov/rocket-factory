package converter

import (
	"github.com/andredubov/rocket-factory/order/internal/model"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

// OrderToPayOrderRequest
func OrderToPayOrderRequest(order *model.Order) *payment_v1.PayOrderRequest {
	return &payment_v1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: convertModelPaymentMethodToProto(order.PaymentInfo.PaymentMethod),
	}
}

// convertModelPaymentMethodToProto конвертирует model.PaymentMethod в payment_v1.PaymentMethod.
func convertModelPaymentMethodToProto(method model.PaymentMethod) payment_v1.PaymentMethod {
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
