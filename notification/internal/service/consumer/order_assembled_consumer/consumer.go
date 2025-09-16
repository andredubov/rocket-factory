package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/andredubov/rocket-factory/notification/internal/converter/kafka"
	"github.com/andredubov/rocket-factory/notification/internal/service"
	"github.com/andredubov/rocket-factory/platform/pkg/kafka"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

type consumerService struct {
	orderAssembledEventConsumer kafka.Consumer
	orderAssembledEventDecoder  kafkaConverter.OrderAssembledEventDecoder
	telegramService             service.TelegramService
}

func NewService(
	orderAssembledConsumer kafka.Consumer,
	orderAssembledDecoder kafkaConverter.OrderAssembledEventDecoder,
	telegramService service.TelegramService,
) *consumerService {
	return &consumerService{
		orderAssembledEventConsumer: orderAssembledConsumer,
		orderAssembledEventDecoder:  orderAssembledDecoder,
		telegramService:             telegramService,
	}
}

func (s *consumerService) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 starting OrderAssembledEventConsumer service")

	err := s.orderAssembledEventConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "❌ consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}
