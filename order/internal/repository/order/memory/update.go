package memory

import (
	"context"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
)

// UpdateOrder modifies an existing order in the repository.
func (r *ordersRepository) UpdateOrder(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.OrderUUID]; !exists {
		return repository.ErrOrderNotFoundWith(order.OrderUUID)
	}

	// Store a copy of the order to prevent external modifications
	orderCopy := converter.OrderToRepoModel(order)
	r.orders[order.OrderUUID] = &orderCopy
	return nil
}
