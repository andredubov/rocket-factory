package memory

import (
	"context"

	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// UpdateOrder modifies an existing order in the repository.
func (r *ordersRepository) UpdateOrder(ctx context.Context, order model.Order) error {
	if !order.Status.IsValid() {
		return repository.ErrInvalidOrderStatusWith(order.Status)
	}
	if order.PaymentInfo != nil && !order.PaymentInfo.PaymentMethod.IsValid() {
		return repository.ErrInvalidPaymentMethodWith(order.PaymentInfo.PaymentMethod)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.OrderUUID]; !exists {
		return repository.ErrOrderNotFoundWith(order.OrderUUID)
	}

	// Store a copy of the order to prevent external modifications
	orderCopy := order
	r.orders[order.OrderUUID] = &orderCopy
	return nil
}
