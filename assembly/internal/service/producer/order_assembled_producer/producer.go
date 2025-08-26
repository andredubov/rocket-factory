package order_assembled_producer

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/assembly/internal/converter"
	"github.com/andredubov/rocket-factory/assembly/internal/model"
	"github.com/andredubov/rocket-factory/platform/pkg/kafka"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

type service struct {
	orderAssembledProducer kafka.Producer
}

func NewService(orderAssembledProducer kafka.Producer) *service {
	return &service{
		orderAssembledProducer: orderAssembledProducer,
	}
}

func (p *service) ProduceOrderAssembledEvent(ctx context.Context, event model.OrderAssembledEvent) error {
	msg := converter.OrderAssembledEventToProtobufEvent(event)

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderAssembledEvent", zap.Error(err))
		return err
	}

	err = p.orderAssembledProducer.Send(ctx, []byte(event.UUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderAssembledEvent", zap.Error(err))
		return err
	}

	return nil
}
