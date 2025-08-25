package order

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/andredubov/rocket-factory/order/internal/converter/kafka"
	"github.com/andredubov/rocket-factory/order/internal/service"
	"github.com/andredubov/rocket-factory/platform/pkg/kafka"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

type consumerService struct {
	orderAssembledEventConsumer kafka.Consumer
	orderAssembledEventDecoder  kafkaConverter.OrderAssembledEventDecoder
	ordersRepository            service.OrdersRepository
}

func NewService(
	orderAssembledConsumer kafka.Consumer,
	orderAssembledDecoder kafkaConverter.OrderAssembledEventDecoder,
	ordersRepository service.OrdersRepository,
) *consumerService {
	return &consumerService{
		orderAssembledEventConsumer: orderAssembledConsumer,
		orderAssembledEventDecoder:  orderAssembledDecoder,
		ordersRepository:            ordersRepository,
	}
}

func (s *consumerService) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting OrderAssembledEventConsumer service")

	err := s.orderAssembledEventConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}
