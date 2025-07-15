package orders

import (
	"context"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
)

// UpdateOrder modifies an existing order.
func (s *ordersService) UpdateOrder(ctx context.Context, order model.Order) error {
	repoOrder := converter.OrderToRepoModel(order)
	return s.ordersRepository.UpdateOrder(ctx, repoOrder)
}
