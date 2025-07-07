package memory

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/andredubov/rocket-factory/order/internal/repository"
)

// DeleteOrder removes an order from the repository by its UUID.
// Returns an error if the order doesn't exist.
// Thread-safe: uses mutex for synchronization.
func (r *ordersRepository) DeleteOrder(ctx context.Context, uuid uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[uuid]; !exists {
		return fmt.Errorf("%w: %s", repository.ErrOrderNotFound, uuid)
	}

	delete(r.orders, uuid)
	return nil
}
