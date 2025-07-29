package handler

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/client/grpc"
	"github.com/andredubov/rocket-factory/order/internal/model"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

// Orders defines the interface for order service operations.
type OrdersService interface {
	GetOrder(ctx context.Context, uuid uuid.UUID) (*model.Order, error)
	AddOrder(ctx context.Context, order model.Order) error
	UpdateOrder(ctx context.Context, order model.Order) error
	DeleteOrder(ctx context.Context, uuid uuid.UUID) error
	GetUserOrders(ctx context.Context, userUUID uuid.UUID) ([]model.Order, error)
}

// OrderImplementation реализует интерфейс обработчика заказов.
type OrderImplementation struct {
	order_v1.UnimplementedHandler
	ordersService   OrdersService
	paymentClient   grpc.PaymentClient
	inventoryClient grpc.InventoryClient
}

// NewOrderHandler создает новый экземпляр обработчика заказов.
func NewOrderHandler(
	service OrdersService,
	paymentClient grpc.PaymentClient,
	inventoryClient grpc.InventoryClient,
) *OrderImplementation {
	return &OrderImplementation{
		ordersService:   service,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
	}
}
