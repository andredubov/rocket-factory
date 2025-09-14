package orders

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/order/internal/converter"
	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	"github.com/andredubov/rocket-factory/platform/pkg/tracing"
)

func (s *ordersService) PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod string) (*model.Order, error) {
	ctx, span := tracing.StartSpan(ctx, "order.call_payment",
		trace.WithAttributes(
			attribute.String("order.id", orderUUID.String()),
		),
	)
	defer span.End()

	order, err := s.ordersRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx, "❌ [OrderService] Не удалось получить заказ", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Валидация статуса заказа
	if order.Status != model.OrderStatusPending {
		logger.Error(ctx, "❌ [OrderService] Заказ не может быть оплачен", zap.String("status", string(order.Status)))
		span.SetAttributes(attribute.String("status", string(order.Status)))
		span.RecordError(model.ErrInvalidOrderStatus)
		span.SetStatus(codes.Error, model.ErrInvalidOrderStatus.Error())
		return nil, model.ErrInvalidOrderStatus
	}

	// Валидация метода оплаты заказа
	payment, err := model.NewPaymentMethod(paymentMethod)
	if err != nil {
		logger.Error(ctx, "❌ [OrderService] failed to vaild payment method", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Подготовка платежной информации
	order.PaymentInfo = &model.PaymentInfo{
		PaymentMethod: payment,
	}

	// Вызов платежного сервиса
	transactionUUID, err := s.paymentClient.PayOrder(ctx, order)
	if err != nil {
		logger.Error(ctx, "❌ [OrderService] failed to pay order", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Обновление информации о заказе
	order.PaymentInfo.TransactionUUID = transactionUUID
	order.Status = model.OrderStatusPaid

	updateInfo := converter.OrderToOrderUpdateInfo(order)

	// Сохранение обновленного заказа
	if err := s.ordersRepository.UpdateOrder(ctx, updateInfo); err != nil {
		logger.Error(ctx, "❌ [OrderService] failed to updated order status", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	event := converter.OrderToOrderPaidEvent(order)

	// Отправка сообщения об успешной оплате заказа в kafka
	if err := s.orderPaidEventProducer.ProduceOrderPaidEvent(ctx, event); err != nil {
		logger.Error(ctx, "❌ [OrderService] failed to produce OrderPaid event", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(
		attribute.String(
			"order.transaction_uuid",
			order.PaymentInfo.TransactionUUID.String(),
		),
	)
	span.SetStatus(codes.Ok, "order payment succeeded")

	return order, nil
}
