package order

import (
	"context"

	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

func (i *Implementation) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
	return nil, nil
}
