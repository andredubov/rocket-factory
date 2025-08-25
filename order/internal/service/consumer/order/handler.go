package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/order/internal/converter"
	"github.com/andredubov/rocket-factory/platform/pkg/kafka"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

func (c *consumerService) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := c.orderAssembledEventDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembledEvent", zap.Error(err))
		return err
	}

	updateInfo, err := converter.OrderAssembledEventToOrderUpdateInfo(event)
	if err != nil {
		logger.Error(ctx, "Failed to convert OrderAssembledEvent to OrderUpdateInfo", zap.Error(err))
		return err
	}

	if err := c.ordersRepository.UpdateOrder(ctx, updateInfo); err != nil {
		logger.Error(ctx, "Failed to update order status on AssembledStatus", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.UUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	return nil
}
