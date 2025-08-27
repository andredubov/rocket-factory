package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/platform/pkg/kafka"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

func (c *consumerService) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := c.orderPaidEventDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaidEvent", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing OrderPaidEvent",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.UUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.String("transaction_uuid", event.TransactionUUID),
		zap.String("payment_method", string(event.PaymentMethod)),
	)

	err = c.telegramService.SendOrderPaidNotification(ctx, event.UUID, event)
	if err != nil {
		logger.Warn(ctx, "Notification of paid order failed to send", zap.Error(err))
		return err
	}

	return nil
}
