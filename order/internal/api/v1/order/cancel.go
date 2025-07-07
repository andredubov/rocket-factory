package order

import (
	"context"

	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

func (i *Implementation) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	return nil, nil
}
