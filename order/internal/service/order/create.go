package orders

import (
	"context"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// AddOrder creates a new order.
func (s *ordersService) AddOrder(ctx context.Context, order model.Order) error {
	return s.ordersRepository.AddOrder(ctx, order)
}
