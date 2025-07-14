package orders

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

func (s *ordersService) GetOrder(ctx context.Context, uuid uuid.UUID) (*model.Order, error) {
	return nil, nil
}

func (s *ordersService) GetUserOrders(ctx context.Context, userUUID uuid.UUID) ([]model.Order, error) {
	return nil, nil
}
