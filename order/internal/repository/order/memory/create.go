package memory

import (
	"context"

	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// AddOrder adds a new order to the repository.
// Validates order status and payment method before adding.
// Returns an error if:
//   - order status is invalid
//   - payment method is invalid
//   - order with the same UUID already exists
//
// Thread-safe: uses mutex for synchronization.
func (r *ordersRepository) AddOrder(ctx context.Context, order model.Order) error {
	if !order.Status.IsValid() {
		return repository.ErrInvalidOrderStatusWith(order.Status)
	}
	if order.PaymentInfo != nil && !order.PaymentInfo.PaymentMethod.IsValid() {
		return repository.ErrInvalidPaymentMethodWith(order.PaymentInfo.PaymentMethod)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.OrderUUID]; exists {
		return repository.ErrOrderAlreadyExistsWith(order.OrderUUID)
	}

	// Store a copy of the order to prevent external modifications
	orderCopy := order
	r.orders[order.OrderUUID] = &orderCopy
	return nil
}
