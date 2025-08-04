package orders

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// GetOrder retrieves an order by its UUID.
func (s *ordersService) GetOrder(ctx context.Context, uuid uuid.UUID) (*model.Order, error) {
	order, err := s.ordersRepository.GetOrder(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return order, nil
}
