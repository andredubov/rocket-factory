package order

import (
	"context"

	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

func (i *Implementation) GetOrderByUuid(ctx context.Context, params order_v1.GetOrderByUuidParams) (order_v1.GetOrderByUuidRes, error) {
	return nil, nil
}
