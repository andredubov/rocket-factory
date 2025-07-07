package repository

import (
	"context"
	"fmt"

	"github.com/andredubov/rocket-factory/order/internal/repository/model"
	"github.com/gofrs/uuid"
)

// Error definitions
var (
	ErrInvalidOrderStatus   = fmt.Errorf("invalid order status")
	ErrInvalidPaymentMethod = fmt.Errorf("invalid payment method")
	ErrOrderAlreadyExists   = fmt.Errorf("order already exists")
	ErrOrderNotFound        = fmt.Errorf("order not found")
)

// Orders defines the interface for order repository operations.
// All implementations must be thread-safe.
type Orders interface {
	// GetOrder retrieves an order by its UUID
	GetOrder(ctx context.Context, uuid uuid.UUID) (*model.Order, error)

	// AddOrder creates a new order
	AddOrder(ctx context.Context, order model.Order) error

	// UpdateOrder modifies an existing order
	UpdateOrder(ctx context.Context, order model.Order) error

	// DeleteOrder removes an order by its UUID
	DeleteOrder(ctx context.Context, uuid uuid.UUID) error

	// GetUserOrders retrieves all orders for a specific user
	GetUserOrders(ctx context.Context, userUUID uuid.UUID) ([]model.Order, error)
}
