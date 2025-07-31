package handler

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/service"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

// Orders defines the interface for order service operations.
type OrdersService interface {
	GetOrder(ctx context.Context, uuid uuid.UUID) (*model.Order, error)
	CreateOrder(ctx context.Context, order model.Order) error
	CancelOrder(ctx context.Context, uuid uuid.UUID) error
	PayOrder(ctx context.Context, uuid uuid.UUID, paymentMethod string) (*model.Order, error)
}

// OrderImplementation реализует интерфейс обработчика заказов.
type OrderImplementation struct {
	order_v1.UnimplementedHandler
	ordersService   OrdersService
	paymentClient   service.PaymentClient
	inventoryClient service.InventoryClient
}

// NewOrderHandler создает новый экземпляр обработчика заказов.
func NewOrderHandler(
	service OrdersService,
	paymentClient service.PaymentClient,
	inventoryClient service.InventoryClient,
) *OrderImplementation {
	return &OrderImplementation{
		ordersService:   service,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
	}
}
