package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// Error definitions
var (
	ErrInvalidOrderStatus   = errors.New("invalid order status")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrOrderAlreadyExists   = errors.New("order already exists")
	ErrOrderNotFound        = errors.New("order not found")
)

// Helper functions for creating formatted errors
func ErrInvalidOrderStatusWith(status model.OrderStatus) error {
	return fmt.Errorf("%w: %s", ErrInvalidOrderStatus, status)
}

func ErrInvalidPaymentMethodWith(method model.PaymentMethod) error {
	return fmt.Errorf("%w: %s", ErrInvalidPaymentMethod, method)
}

func ErrOrderAlreadyExistsWith(uuid uuid.UUID) error {
	return fmt.Errorf("%w: %s", ErrOrderAlreadyExists, uuid)
}

func ErrOrderNotFoundWith(uuid uuid.UUID) error {
	return fmt.Errorf("%w: %s", ErrOrderNotFound, uuid)
}

// Orders defines the interface for order repository operations.
type Orders interface {
	GetOrder(ctx context.Context, uuid uuid.UUID) (*model.Order, error)
	AddOrder(ctx context.Context, order model.Order) error
	UpdateOrder(ctx context.Context, order model.Order) error
	DeleteOrder(ctx context.Context, uuid uuid.UUID) error
	GetUserOrders(ctx context.Context, userUUID uuid.UUID) ([]model.Order, error)
}
