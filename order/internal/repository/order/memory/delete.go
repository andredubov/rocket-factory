package memory

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/repository"
)

// DeleteOrder removes an order from the repository by its UUID.
func (r *ordersRepository) DeleteOrder(ctx context.Context, uuid uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[uuid]; !exists {
		return repository.ErrOrderNotFoundWith(uuid)
	}

	delete(r.orders, uuid)
	return nil
}
