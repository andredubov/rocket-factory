package handler

import (
	"github.com/andredubov/rocket-factory/order/internal/client/grpc"
	"github.com/andredubov/rocket-factory/order/internal/service"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

// OrderImplementation реализует интерфейс обработчика заказов.
type OrderImplementation struct {
	order_v1.UnimplementedHandler
	ordersService   service.Orders
	paymentClient   grpc.PaymentClient
	inventoryClient grpc.InventoryClient
}

// NewOrderHandler создает новый экземпляр обработчика заказов.
func NewOrderHandler(
	service service.Orders,
	paymentClient grpc.PaymentClient,
	inventoryClient grpc.InventoryClient,
) *OrderImplementation {
	return &OrderImplementation{
		ordersService:   service,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
	}
}
