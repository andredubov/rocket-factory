package grpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// type InventoryClient interface {
// 	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
// }

type PaymentClient interface {
	PayOrder(ctx context.Context, order *model.Order) (uuid.UUID, error)
}
