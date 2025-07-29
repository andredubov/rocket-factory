package orders

import (
	"context"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// UpdateOrder modifies an existing order.
func (s *ordersService) UpdateOrder(ctx context.Context, order model.Order) error {
	return s.ordersRepository.UpdateOrder(ctx, order)
}
