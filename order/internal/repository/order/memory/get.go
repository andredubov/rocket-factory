package memory

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// GetUserOrders retrieves all orders belonging to a specific user.
// Returns an empty slice if no orders found for the user.
// Thread-safe: uses RWMutex for concurrent reads.
func (r *ordersRepository) GetUserOrders(ctx context.Context, userUUID uuid.UUID) ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userOrders []model.Order
	for _, order := range r.orders {
		if order.UserUUID == userUUID {
			// Add a copy to prevent external modifications
			orderCopy := *order
			userOrders = append(userOrders, orderCopy)
		}
	}

	return userOrders, nil
}

// GetOrder retrieves a single order by its UUID.
// Returns an error if the order doesn't exist.
// Thread-safe: uses RWMutex for concurrent reads.
func (r *ordersRepository) GetOrder(ctx context.Context, uuid uuid.UUID) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[uuid]
	if !exists {
		return nil, repository.ErrOrderNotFoundWith(uuid)
	}

	// Return a copy to prevent external modifications
	orderCopy := *order
	return &orderCopy, nil
}
