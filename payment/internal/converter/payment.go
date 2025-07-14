package converter

import (
	"github.com/andredubov/rocket-factory/payment/internal/model"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

// PaymentFromRequest converts a gRPC PayOrderRequest into the domain Payment model.
func PaymentFromRequest(r *payment_v1.PayOrderRequest) model.Payment {
	return model.Payment{
		UserID:        r.GetUserUuid(),
		OrderID:       r.GetOrderUuid(),
		PaymentMethod: model.PaymentMethod(r.GetPaymentMethod()),
	}
}
