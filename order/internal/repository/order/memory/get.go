package memory

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
)

// GetUserOrders retrieves all orders belonging to a specific user.
func (r *ordersRepository) GetUserOrders(ctx context.Context, userUUID uuid.UUID) ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userOrders := make([]model.Order, 0)
	for _, repoOrder := range r.orders {
		if repoOrder.UserUUID == userUUID {
			// Add a copy to prevent external modifications
			order := converter.OrderToModel(*repoOrder)
			orderCopy := *order
			userOrders = append(userOrders, orderCopy)
		}
	}

	return userOrders, nil
}

// GetOrder retrieves a single order by its UUID.
func (r *ordersRepository) GetOrder(ctx context.Context, uuid uuid.UUID) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoOrder, exists := r.orders[uuid]
	if !exists {
		return nil, repository.ErrOrderNotFoundWith(uuid)
	}

	// Return a copy to prevent external modifications
	order := converter.OrderToModel(*repoOrder)
	orderCopy := order
	return orderCopy, nil
}
