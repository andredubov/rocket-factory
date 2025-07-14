package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// Orders defines the interface for order service operations.
type Orders interface {
	GetOrder(ctx context.Context, uuid uuid.UUID) (*model.Order, error)
	AddOrder(ctx context.Context, order model.Order) error
	UpdateOrder(ctx context.Context, order model.Order) error
	DeleteOrder(ctx context.Context, uuid uuid.UUID) error
	GetUserOrders(ctx context.Context, userUUID uuid.UUID) ([]model.Order, error)
}
