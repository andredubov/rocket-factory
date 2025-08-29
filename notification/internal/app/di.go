package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	httpClient "github.com/andredubov/rocket-factory/notification/internal/client/http"
	telegramClient "github.com/andredubov/rocket-factory/notification/internal/client/http/telegram"
	"github.com/andredubov/rocket-factory/notification/internal/config"
	kafkaConverter "github.com/andredubov/rocket-factory/notification/internal/converter/kafka"
	"github.com/andredubov/rocket-factory/notification/internal/converter/kafka/decoder"
	"github.com/andredubov/rocket-factory/notification/internal/service"
	"github.com/andredubov/rocket-factory/notification/internal/service/consumer/order_assembled_consumer"
	"github.com/andredubov/rocket-factory/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/andredubov/rocket-factory/notification/internal/service/telegram"
	"github.com/andredubov/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/andredubov/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/andredubov/rocket-factory/platform/pkg/kafka/consumer"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/andredubov/rocket-factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	orderPaidEventConsumerGroup        sarama.ConsumerGroup
	orderAssembledEventConsumerGroup   sarama.ConsumerGroup
	orderPaidEventConsumer             wrappedKafka.Consumer
	orderAssembledEventConsumer        wrappedKafka.Consumer
	orderPaidEventDecoder              kafkaConverter.OrderPaidEventDecoder
	orderAssembledEventDecoder         kafkaConverter.OrderAssembledEventDecoder
	orderPaidEventConsumerService      service.ConsumerService
	orderAssembledEventConsumerService service.ConsumerService

	telegramClient  httpClient.TelegramClient
	telegramBot     *bot.Bot
	telegramService service.TelegramService
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().Telegram.TelegramBotToken())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}

func (d *diContainer) TelegramClient(ctx context.Context) httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(
			d.TelegramClient(ctx),
		)
	}

	return d.telegramService
}

func (d *diContainer) OrderPaidEventConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidEventConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidEventConsumer.GroupID(),
			config.AppConfig().OrderPaidEventConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create OrderPaidEvent consumer group: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.orderPaidEventConsumerGroup.Close()
		})

		d.orderPaidEventConsumerGroup = consumerGroup
	}

	return d.orderPaidEventConsumerGroup
}

func (d *diContainer) OrderAssebledEventConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledEventConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledEventConsumer.GroupID(),
			config.AppConfig().OrderAssembledEventConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create OrderAssembledEvent consumer group: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.orderAssembledEventConsumerGroup.Close()
		})

		d.orderAssembledEventConsumerGroup = consumerGroup
	}

	return d.orderAssembledEventConsumerGroup
}

func (d *diContainer) OrderPaidEventConsumer() wrappedKafka.Consumer {
	if d.orderPaidEventConsumer == nil {
		d.orderPaidEventConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidEventConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidEventConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidEventConsumer
}

func (d *diContainer) OrderAssembledEventConsumer() wrappedKafka.Consumer {
	if d.orderAssembledEventConsumer == nil {
		d.orderAssembledEventConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidEventConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledEventConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledEventConsumer
}

func (d *diContainer) OrderPaidEventDecoder() kafkaConverter.OrderPaidEventDecoder {
	if d.orderPaidEventDecoder == nil {
		d.orderPaidEventDecoder = decoder.NewOrderPaidEventDecoder()
	}

	return d.orderPaidEventDecoder
}

func (d *diContainer) OrderAssembledEventDecoder() kafkaConverter.OrderAssembledEventDecoder {
	if d.orderAssembledEventDecoder == nil {
		d.orderAssembledEventDecoder = decoder.NewOrderAssembledEventDecoder()
	}

	return d.orderAssembledEventDecoder
}

func (d *diContainer) OrderPaidEventConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderPaidEventConsumerService == nil {
		d.orderPaidEventConsumerService = order_paid_consumer.NewService(
			d.OrderPaidEventConsumer(),
			d.OrderPaidEventDecoder(),
			d.TelegramService(ctx),
		)
	}

	return d.orderPaidEventConsumerService
}

func (d *diContainer) OrderAssembledEventConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderAssembledEventConsumerService == nil {
		d.orderAssembledEventConsumerService = order_assembled_consumer.NewService(
			d.OrderAssembledEventConsumer(),
			d.OrderAssembledEventDecoder(),
			d.TelegramService(ctx),
		)
	}

	return d.orderAssembledEventConsumerService
}
