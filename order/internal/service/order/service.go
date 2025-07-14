package orders

import (
	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/service"
)

type ordersService struct {
	ordersRepository repository.Orders
}

func NewService(repo repository.Orders) service.Orders {
	return &ordersService{
		ordersRepository: repo,
	}
}
