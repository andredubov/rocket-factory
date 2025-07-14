package handler

import (
	"github.com/andredubov/rocket-factory/order/internal/service"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

// OrderImplementation реализует интерфейс обработчика заказов.
type OrderImplementation struct {
	order_v1.UnimplementedHandler
	ordersService   service.Orders
	paymentClient   payment_v1.PaymentServiceClient
	inventoryClient inventory_v1.InventoryServiceClient
}

// NewOrderHandler создает новый экземпляр обработчика заказов.
func NewOrderHandler(
	service service.Orders,
	paymentClient payment_v1.PaymentServiceClient,
	inventoryClient inventory_v1.InventoryServiceClient,
) *OrderImplementation {
	return &OrderImplementation{
		ordersService:   service,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
	}
}
