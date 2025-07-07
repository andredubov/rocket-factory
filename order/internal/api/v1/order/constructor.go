package order

import (
	"github.com/andredubov/rocket-factory/order/internal/repository"
	order_v1 "github.com/andredubov/rocket-factory/shared/pkg/openapi/order/v1"
)

type Implementation struct {
	order_v1.UnimplementedHandler
	ordersRepository repository.Orders
}

func NewOrderHandler(repo repository.Orders) *Implementation {
	return &Implementation{
		ordersRepository: repo,
	}
}
