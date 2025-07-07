package memory

import (
	"context"
	"fmt"

	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// UpdateOrder modifies an existing order in the repository.
// Validates order status and payment method before updating.
// Returns an error if:
//   - order status is invalid
//   - payment method is invalid
//   - order doesn't exist
//
// Thread-safe: uses mutex for synchronization.
func (r *ordersRepository) UpdateOrder(ctx context.Context, order model.Order) error {
	if !order.Status.IsValid() {
		return fmt.Errorf("%w: %s", repository.ErrInvalidOrderStatus, order.Status)
	}
	if order.PaymentInfo != nil && !order.PaymentInfo.PaymentMethod.IsValid() {
		return fmt.Errorf("%w: %s", repository.ErrInvalidPaymentMethod, order.PaymentInfo.PaymentMethod)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.OrderUUID]; !exists {
		return fmt.Errorf("%w: %s", repository.ErrOrderNotFound, order.OrderUUID)
	}

	// Store a copy of the order to prevent external modifications
	orderCopy := order
	r.orders[order.OrderUUID] = &orderCopy
	return nil
}
