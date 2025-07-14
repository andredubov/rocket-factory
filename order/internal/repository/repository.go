package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// Helper functions for creating formatted errors
func ErrInvalidOrderStatusWith(status repoModel.OrderStatus) error {
	return fmt.Errorf("%w: %s", model.ErrInvalidOrderStatus, status)
}

func ErrInvalidPaymentMethodWith(method repoModel.PaymentMethod) error {
	return fmt.Errorf("%w: %s", model.ErrInvalidPaymentMethod, method)
}

func ErrOrderAlreadyExistsWith(uuid uuid.UUID) error {
	return fmt.Errorf("%w: %s", model.ErrOrderAlreadyExists, uuid)
}

func ErrOrderNotFoundWith(uuid uuid.UUID) error {
	return fmt.Errorf("%w: %s", model.ErrOrderNotFound, uuid)
}

// Orders defines the interface for order repository operations.
type Orders interface {
	GetOrder(ctx context.Context, uuid uuid.UUID) (*repoModel.Order, error)
	AddOrder(ctx context.Context, order repoModel.Order) error
	UpdateOrder(ctx context.Context, order repoModel.Order) error
	DeleteOrder(ctx context.Context, uuid uuid.UUID) error
	GetUserOrders(ctx context.Context, userUUID uuid.UUID) ([]repoModel.Order, error)
}
