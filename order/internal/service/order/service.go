package orders

import (
	handler "github.com/andredubov/rocket-factory/order/internal/api/v1/order"
	"github.com/andredubov/rocket-factory/order/internal/service"
)

// ordersService implements the service.Orders interface.
type ordersService struct {
	ordersRepository       service.OrdersRepository
	paymentClient          service.PaymentClient
	inventoryClient        service.InventoryClient
	orderPaidEventProducer service.ProducerService
}

// NewService creates a new instance of the order service.
func NewService(
	repository service.OrdersRepository,
	paymentClient service.PaymentClient,
	inventoryClient service.InventoryClient,
	orderPaidEventProducer service.ProducerService,
) handler.OrdersService {
	return &ordersService{
		ordersRepository:       repository,
		paymentClient:          paymentClient,
		inventoryClient:        inventoryClient,
		orderPaidEventProducer: orderPaidEventProducer,
	}
}
