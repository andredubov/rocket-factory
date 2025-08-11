package memory

import (
	"context"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
)

// AddOrder adds a new order to the repository.
func (r *ordersRepository) AddOrder(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	repoOrder := converter.OrderToRepoModel(order)
	if _, exists := r.orders[repoOrder.OrderUUID]; exists {
		return model.ErrOrderAlreadyExistsWith(order.OrderUUID)
	}

	// Store a copy of the order to prevent external modifications
	orderCopy := repoOrder
	r.orders[order.OrderUUID] = &orderCopy
	return nil
}
