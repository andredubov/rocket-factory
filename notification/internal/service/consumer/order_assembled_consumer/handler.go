package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/platform/pkg/kafka"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

func (c *consumerService) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := c.orderAssembledEventDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembledEvent", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing OrderAssembledEvent",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.UUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int64("build__time_sec", event.BuildTimeSec),
	)

	err = c.telegramService.SendOrderAssembledNotification(ctx, event.UUID, event)
	if err != nil {
		logger.Warn(ctx, "Notification of assembled order failed to send", zap.Error(err))
		return err
	}

	return nil
}
