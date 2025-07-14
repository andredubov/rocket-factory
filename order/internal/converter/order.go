package converter

import (
	"github.com/andredubov/rocket-factory/order/internal/model"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

func OrderFromCreateOrderRequest(req *order_v1.CreateOrderRequest) model.Order {
	return model.Order{
		UserUUID:  req.GetUserUUID(),
		PartUUIDs: req.GetPartUuids(),
	}
}

func OrderToCreateOrderResponse(order model.Order) *order_v1.CreateOrderResponse {
	return &order_v1.CreateOrderResponse{
		OrderUUID:  order_v1.NewOptUUID(order.OrderUUID),
		TotalPrice: order_v1.NewOptFloat64(order.TotalPrice),
	}
}
