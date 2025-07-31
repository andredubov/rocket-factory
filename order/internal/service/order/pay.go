package orders

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

func (s *ordersService) PayOrder(ctx context.Context, uuid uuid.UUID, paymentMethod string) (*model.Order, error) {
	order, err := s.ordersRepository.GetOrder(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// Валидация статуса заказа
	if order.Status != model.OrderStatusPending {
		return nil, model.ErrInvalidOrderStatus
	}

	// Валидация метода оплаты заказа
	payment, err := model.NewPaymentMethod(paymentMethod)
	if err != nil {
		return nil, err
	}

	// Подготовка платежной информации
	order.PaymentInfo = &model.PaymentInfo{
		PaymentMethod: payment,
	}

	// Вызов платежного сервиса
	transactionUUID, err := s.paymentClient.PayOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// Обновление информации о заказе
	order.PaymentInfo.TransactionUUID = transactionUUID
	order.Status = model.OrderStatusPaid

	// Сохранение обновленного заказа
	if err := s.ordersRepository.UpdateOrder(ctx, *order); err != nil {
		return nil, err
	}

	return order, nil
}
