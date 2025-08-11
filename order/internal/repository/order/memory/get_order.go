package memory

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
)

// GetOrder retrieves a single order by its UUID.
func (r *ordersRepository) GetOrder(ctx context.Context, uuid uuid.UUID) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoOrder, exists := r.orders[uuid]
	if !exists {
		return nil, model.ErrOrderNotFoundWith(uuid)
	}

	// Return a copy to prevent external modifications
	order := converter.OrderToModel(*repoOrder)
	orderCopy := order
	return orderCopy, nil
}
