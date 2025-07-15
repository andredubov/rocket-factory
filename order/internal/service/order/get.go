package orders

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
)

// GetOrder retrieves an order by its UUID.
func (s *ordersService) GetOrder(ctx context.Context, uuid uuid.UUID) (*model.Order, error) {
	order, err := s.ordersRepository.GetOrder(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return converter.OrderToModel(*order), nil
}

// GetUserOrders returns all orders belonging to a user by their UUID.
func (s *ordersService) GetUserOrders(ctx context.Context, userUUID uuid.UUID) ([]model.Order, error) {
	orders, err := s.ordersRepository.GetUserOrders(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return converter.OrdersToModel(orders), nil
}
