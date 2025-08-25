package order

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/assembly/internal/model"
	"github.com/andredubov/rocket-factory/platform/pkg/kafka"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

func (c *consumerService) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := c.orderPaidEventDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaidEvent", zap.Error(err))
		return err
	}

	go func(ctx context.Context) {
		var (
			delay               = 10 * time.Second
			begin               = time.Now()
			orderAssembledEvent = model.OrderAssembledEvent{
				UUID:         event.UUID,
				OrderUUID:    event.OrderUUID,
				UserUUID:     event.UserUUID,
				BuildTimeSec: int64(time.Since(begin)),
			}
		)

		logger.Info(ctx, "Processing OrderPaidEvent",
			zap.String("topic", msg.Topic),
			zap.Any("partition", msg.Partition),
			zap.Any("offset", msg.Offset),
			zap.String("event_uuid", event.UUID),
			zap.String("order_uuid", event.OrderUUID),
			zap.String("user_uuid", event.UserUUID),
		)

		select {
		case <-time.After(delay):
			err = c.orderAssembledProducer.ProduceOrderAssembledEvent(ctx, orderAssembledEvent)
			if err != nil {
				logger.Info(ctx, "Failed to produce OrderAssembledEvent", zap.Error(err))
			}
		case <-ctx.Done():
		}
	}(ctx)

	return nil
}
