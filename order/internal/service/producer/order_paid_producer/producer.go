package order_paid_producer

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/order/internal/converter"
	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/platform/pkg/kafka"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

type service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (p *service) ProduceOrderPaidEvent(ctx context.Context, event model.OrderPaidEvent) error {
	msg := converter.OrderPaidEventToProtobufEvent(event)

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderPaidEvent", zap.Error(err))
		return err
	}

	err = p.orderPaidProducer.Send(ctx, []byte(event.UUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderPaidEvent", zap.Error(err))
		return err
	}

	return nil
}
