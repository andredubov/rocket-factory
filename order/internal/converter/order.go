package converter

import (
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
