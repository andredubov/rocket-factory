package orders

import (
	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/service"
)

// ordersService implements the service.Orders interface.
type ordersService struct {
	ordersRepository repository.Orders
}

// NewService creates a new instance of the order service.
// Accepts an orders repository and returns an implementation of service.Orders interface
func NewService(repo repository.Orders) service.Orders {
	return &ordersService{
		ordersRepository: repo,
	}
}
