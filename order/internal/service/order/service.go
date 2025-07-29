package orders

import (
	"github.com/andredubov/rocket-factory/order/internal/service"
)

// ordersService implements the service.Orders interface.
type ordersService struct {
	ordersRepository service.OrdersRepository
}

// NewService creates a new instance of the order service.
func NewService(repo service.OrdersRepository) service.OrdersRepository {
	return &ordersService{
		ordersRepository: repo,
	}
}
