package order

import (
	"github.com/andredubov/rocket-factory/order/internal/repository"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

type Implementation struct {
	order_v1.UnimplementedHandler
	ordersRepository repository.Orders
	paymentClient    payment_v1.PaymentServiceClient
	inventoryClient  inventory_v1.InventoryServiceClient
}

func NewOrderHandler(
	repo repository.Orders,
	paymentClient payment_v1.PaymentServiceClient,
	inventoryClient inventory_v1.InventoryServiceClient,
) *Implementation {
	return &Implementation{
		ordersRepository: repo,
		paymentClient:    paymentClient,
		inventoryClient:  inventoryClient,
	}
}
