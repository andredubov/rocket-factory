package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/andredubov/rocket-factory/assembly/internal/config"
	kafkaConverter "github.com/andredubov/rocket-factory/assembly/internal/converter/kafka"
	"github.com/andredubov/rocket-factory/assembly/internal/converter/kafka/decoder"
	"github.com/andredubov/rocket-factory/assembly/internal/service"
	"github.com/andredubov/rocket-factory/assembly/internal/service/consumer/order_paid_consumer"
	"github.com/andredubov/rocket-factory/assembly/internal/service/producer/order_assembled_producer"
	"github.com/andredubov/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/andredubov/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/andredubov/rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/andredubov/rocket-factory/platform/pkg/kafka/producer"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/andredubov/rocket-factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	consumerGroup          sarama.ConsumerGroup
	orderPaidEventConsumer wrappedKafka.Consumer
	orderPaidEventDecoder  kafkaConverter.OrderPaidEventDecoder
	consumerService        service.ConsumerService

	syncProducer                sarama.SyncProducer
	orderAssembledEventProducer wrappedKafka.Producer
	producerService             service.ProducerService
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidEventConsumer.GroupID(),
			config.AppConfig().OrderPaidEventConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderPaidEventConsumer() wrappedKafka.Consumer {
	if d.orderPaidEventConsumer == nil {
		d.orderPaidEventConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidEventConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidEventConsumer
}

func (d *diContainer) OrderAssembledEventDecoder() kafkaConverter.OrderPaidEventDecoder {
	if d.orderPaidEventDecoder == nil {
		d.orderPaidEventDecoder = decoder.NewOrderPaidEventDecoder()
	}

	return d.orderPaidEventDecoder
}

func (d *diContainer) ConsumerService(_ context.Context) service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = order_paid_consumer.NewService(
			d.OrderPaidEventConsumer(),
			d.OrderAssembledEventDecoder(),
			d.ProducerService(),
		)
	}

	return d.consumerService
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledEventProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderAssembledEventProducer() wrappedKafka.Producer {
	if d.orderAssembledEventProducer == nil {
		d.orderAssembledEventProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledEventProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderAssembledEventProducer
}

func (d *diContainer) ProducerService() service.ProducerService {
	if d.producerService == nil {
		d.producerService = order_assembled_producer.NewService(d.OrderAssembledEventProducer())
	}

	return d.producerService
}
