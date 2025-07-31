package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// InventoryClient defines the interface for payment service client
type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
}

// PaymentClient defines the interface for payment service client
type PaymentClient interface {
	PayOrder(ctx context.Context, order *model.Order) (uuid.UUID, error)
}

// OrdersRepository defines the interface for order repository operations.
type OrdersRepository interface {
	GetOrder(ctx context.Context, uuid uuid.UUID) (*model.Order, error)
	AddOrder(ctx context.Context, order model.Order) error
	UpdateOrder(ctx context.Context, order model.Order) error
	DeleteOrder(ctx context.Context, uuid uuid.UUID) error
	GetUserOrders(ctx context.Context, userUUID uuid.UUID) ([]model.Order, error)
}
