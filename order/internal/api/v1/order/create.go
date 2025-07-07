package order

import (
	"context"

	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	return nil, nil
}
