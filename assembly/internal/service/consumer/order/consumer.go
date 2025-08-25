package order

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/andredubov/rocket-factory/assembly/internal/converter/kafka"
	"github.com/andredubov/rocket-factory/assembly/internal/service"
	"github.com/andredubov/rocket-factory/platform/pkg/kafka"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

type consumerService struct {
	orderPaidEventConsumer kafka.Consumer
	orderPaidEventDecoder  kafkaConverter.OrderPaidEventDecoder
	orderAssembledProducer service.ProducerService
}

func NewService(
	orderAssembledConsumer kafka.Consumer,
	orderAssembledDecoder kafkaConverter.OrderPaidEventDecoder,
	orderAssembledProducer service.ProducerService,
) *consumerService {
	return &consumerService{
		orderPaidEventConsumer: orderAssembledConsumer,
		orderPaidEventDecoder:  orderAssembledDecoder,
		orderAssembledProducer: orderAssembledProducer,
	}
}

func (s *consumerService) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting OrderPaidEventConsumer service")

	err := s.orderPaidEventConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
