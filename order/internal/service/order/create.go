package orders

import (
	"context"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
)

// AddOrder creates a new order.
func (s *ordersService) AddOrder(ctx context.Context, order model.Order) error {
	repoOrder := converter.OrderToRepoModel(order)
	return s.ordersRepository.AddOrder(ctx, repoOrder)
}
