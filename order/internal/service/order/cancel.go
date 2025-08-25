package orders

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/converter"
	"github.com/andredubov/rocket-factory/order/internal/model"
)

func (s *ordersService) CancelOrder(ctx context.Context, uuid uuid.UUID) error {
	order, err := s.ordersRepository.GetOrder(ctx, uuid)
	if err != nil {
		return err
	}

	// Проверяем статус заказа
	switch order.Status {
	case model.OrderStatusPaid:
		return model.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return model.ErrOrderAlreadyCancelled
	}

	// Меняем статус на Cancelled для заказов в статусе Pending
	order.Status = model.OrderStatusCancelled

	updateInfo := converter.OrderToOrderUpdateInfo(order)

	// Обновляем заказ в репозитории
	return s.ordersRepository.UpdateOrder(ctx, updateInfo)
}
