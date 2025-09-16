package order_paid_consumer

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/assembly/internal/metrics"
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

	logger.Info(ctx, "Processing OrderPaidEvent",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("event_uuid", event.UUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
	)

	return c.produceOrderAssembledEvent(ctx, event)
}

func (c *consumerService) produceOrderAssembledEvent(ctx context.Context, event model.OrderPaidEvent) error {
	begin := time.Now()

	// Имитируем время сборки (10 секунд)
	select {
	case <-time.After(10 * time.Second):
		// Вычисляем реальное время выполнения
		buildTime := time.Since(begin)

		orderAssembledEvent := model.OrderAssembledEvent{
			UUID:         event.UUID,
			OrderUUID:    event.OrderUUID,
			UserUUID:     event.UserUUID,
			BuildTimeSec: int64(buildTime.Seconds()),
		}

		metrics.RecordAssemblyDuration(ctx, buildTime)

		produceCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err := c.orderAssembledProducer.ProduceOrderAssembledEvent(produceCtx, orderAssembledEvent)
		if err != nil {
			logger.Error(ctx, "Failed to produce OrderAssembledEvent",
				zap.Error(err),
				zap.String("order_uuid", event.OrderUUID),
			)

			return err
		}

		logger.Info(ctx, "Successfully produced OrderAssembledEvent",
			zap.String("order_uuid", event.OrderUUID),
			zap.Duration("build_time", buildTime),
		)

	case <-ctx.Done():
		logger.Info(ctx, "Order processing cancelled",
			zap.String("order_uuid", event.OrderUUID),
			zap.String("reason", ctx.Err().Error()),
		)
	}

	return nil
}
