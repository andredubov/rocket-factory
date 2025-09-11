package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/andredubov/rocket-factory/notification/internal/converter/kafka"
	"github.com/andredubov/rocket-factory/notification/internal/service"
	"github.com/andredubov/rocket-factory/platform/pkg/kafka"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

type consumerService struct {
	orderPaidEventConsumer kafka.Consumer
	orderPaidEventDecoder  kafkaConverter.OrderPaidEventDecoder
	telegramService        service.TelegramService
}

func NewService(
	orderPaidEventConsumer kafka.Consumer,
	orderPaidEventDecoder kafkaConverter.OrderPaidEventDecoder,
	telegramService service.TelegramService,
) *consumerService {
	return &consumerService{
		orderPaidEventConsumer: orderPaidEventConsumer,
		orderPaidEventDecoder:  orderPaidEventDecoder,
		telegramService:        telegramService,
	}
}

func (s *consumerService) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 starting OrderPaidEventConsumer service")

	err := s.orderPaidEventConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "❌ consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
