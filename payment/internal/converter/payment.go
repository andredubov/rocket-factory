package converter

import (
	"github.com/andredubov/rocket-factory/payment/internal/model"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

// PaymentFromRequest converts a gRPC PayOrderRequest into the domain Payment model.
func PaymentFromRequest(r *payment_v1.PayOrderRequest) model.Payment {
	return model.Payment{
		UserUuid:      r.GetUserUuid(),
		OrderUuid:     r.GetOrderUuid(),
		PaymentMethod: model.PaymentMethod(r.GetPaymentMethod()),
	}
}

// TransactionUuidToResponse converts a payment transaction ID into grps response
func TransactionUuidToResponse(uuid string) *payment_v1.PayOrderResponse {
	return &payment_v1.PayOrderResponse{
		TransactionUuid: uuid,
	}
}
