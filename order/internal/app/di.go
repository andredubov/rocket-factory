package app

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/andredubov/rocket-factory/order/internal/api/v1/order"
	inventoryClient "github.com/andredubov/rocket-factory/order/internal/client/config/env/inventory"
	paymentClient "github.com/andredubov/rocket-factory/order/internal/client/config/env/payment"
	"github.com/andredubov/rocket-factory/order/internal/client/grpc/inventory/v1"
	"github.com/andredubov/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/andredubov/rocket-factory/order/internal/config"
	"github.com/andredubov/rocket-factory/order/internal/config/env"
	kafkaConverter "github.com/andredubov/rocket-factory/order/internal/converter/kafka"
	"github.com/andredubov/rocket-factory/order/internal/converter/kafka/decoder"
	"github.com/andredubov/rocket-factory/order/internal/repository/order/postgres"
	"github.com/andredubov/rocket-factory/order/internal/service"
	consumer "github.com/andredubov/rocket-factory/order/internal/service/consumer/order"
	orders "github.com/andredubov/rocket-factory/order/internal/service/order"
	"github.com/andredubov/rocket-factory/order/internal/service/producer"
	"github.com/andredubov/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/andredubov/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/andredubov/rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/andredubov/rocket-factory/platform/pkg/kafka/producer"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/andredubov/rocket-factory/platform/pkg/middleware/kafka"
	"github.com/andredubov/rocket-factory/platform/pkg/migrator"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	httpConfig            config.HTTPConfig
	inventoryClientConfig config.GRPCConfig
	paymentClientConfig   config.GRPCConfig
	postgresDBConfig      config.PostgresDBConfig
	inventoryClient       service.InventoryClient
	paymentClient         service.PaymentClient

	consumerGroup               sarama.ConsumerGroup
	orderAssembledEventConsumer wrappedKafka.Consumer
	consumerService             service.ConsumerService

	syncProducer               sarama.SyncProducer
	orderPaidEventProducer     wrappedKafka.Producer
	orderAssembledEventDecoder kafkaConverter.OrderAssembledEventDecoder
	producerService            service.ProducerService

	postgresDB           *pgxpool.Pool
	ordersRepository     service.OrdersRepository
	ordersService        api.OrdersService
	serverImplementation *api.OrderImplementation
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) HTTPConfig() config.HTTPConfig {
	if d.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Printf("failed to get http server config: %s\n", err.Error())
			return nil
		}

		d.httpConfig = cfg
	}

	return d.httpConfig
}

func (d *diContainer) InventoryGRPCConfig() config.GRPCConfig {
	if d.inventoryClientConfig == nil {
		cfg, err := inventoryClient.NewGRPCConfig()
		if err != nil {
			log.Printf("failed to create inventory grpc client config: %v\n", err)
			return nil
		}

		d.inventoryClientConfig = cfg
	}

	return d.inventoryClientConfig
}

func (d *diContainer) PaymentGRPCConfig() config.GRPCConfig {
	if d.paymentClientConfig == nil {
		cfg, err := paymentClient.NewGRPCConfig()
		if err != nil {
			log.Printf("failed to create payment grpc client config: %v\n", err)
			return nil
		}

		d.paymentClientConfig = cfg
	}

	return d.paymentClientConfig
}

func (d *diContainer) InventoryClient() service.InventoryClient {
	if d.inventoryClient == nil {
		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		conn, err := grpc.NewClient(d.InventoryGRPCConfig().Address(), dialOptions...)
		if err != nil {
			log.Printf("failed to create inventory grpc client: %v\n", err)
			return nil
		}

		client := inventory_v1.NewInventoryServiceClient(conn)
		d.inventoryClient = inventory.NewClient(client)
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClient() service.PaymentClient {
	if d.paymentClient == nil {
		dialOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		conn, err := grpc.NewClient(d.PaymentGRPCConfig().Address(), dialOptions...)
		if err != nil {
			log.Printf("failed to create Payment grpc client: %v\n", err)
			return nil
		}

		client := payment_v1.NewPaymentServiceClient(conn)
		d.paymentClient = payment.NewClient(client)
	}

	return d.paymentClient
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledEventConsumer.GroupID(),
			config.AppConfig().OrderAssembledEventConsumer.Config(),
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

func (d *diContainer) OrderAssembledEventConsumer() wrappedKafka.Consumer {
	if d.orderAssembledEventConsumer == nil {
		d.orderAssembledEventConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledEventConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledEventConsumer
}

func (d *diContainer) OrderAssembledEventDecoder() kafkaConverter.OrderAssembledEventDecoder {
	if d.orderAssembledEventDecoder == nil {
		d.orderAssembledEventDecoder = decoder.NewOrderAssembledEventDecoder()
	}

	return d.orderAssembledEventDecoder
}

func (d *diContainer) ConsumerService(ctx context.Context) service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = consumer.NewService(
			d.OrderAssembledEventConsumer(),
			d.OrderAssembledEventDecoder(),
			d.OrdersRepository(ctx),
		)
	}

	return d.consumerService
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidEventProducer.Config(),
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

func (d *diContainer) OrderPaidEventProducer() wrappedKafka.Producer {
	if d.orderPaidEventProducer == nil {
		d.orderPaidEventProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidEventProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderPaidEventProducer
}

func (d *diContainer) ProducerService() service.ProducerService {
	if d.producerService == nil {
		d.producerService = producer.NewService(d.OrderPaidEventProducer())
	}

	return d.producerService
}

func (d *diContainer) PostgresDBConfig() config.PostgresDBConfig {
	if d.postgresDBConfig == nil {
		cfg, err := env.NewPostgresDBConfig()
		if err != nil {
			log.Printf("failed to get Postgres database config: %s", err.Error())
			return nil
		}

		d.postgresDBConfig = cfg
	}

	return d.postgresDBConfig
}

func (d *diContainer) PostgresDatabase(ctx context.Context) *pgxpool.Pool {
	if d.postgresDB == nil {
		dbPool, err := pgxpool.New(ctx, d.PostgresDBConfig().DSN())
		if err != nil {
			log.Printf("failed to connect to database: %v\n", err)
			return nil
		}

		err = dbPool.Ping(ctx)
		if err != nil {
			log.Printf("postgres unawailable: %v\n", err)
			return nil
		}

		d.postgresDB = dbPool

		migratorRunner := migrator.NewMigrator(stdlib.OpenDBFromPool(dbPool), d.PostgresDBConfig().MigrationDirectory())
		err = migratorRunner.Up()
		if err != nil {
			log.Printf("failed to up database migration: %v\n", err)
			return nil
		}
	}

	return d.postgresDB
}

func (d *diContainer) OrdersRepository(ctx context.Context) service.OrdersRepository {
	if d.ordersRepository == nil {
		d.ordersRepository = postgres.NewOrderRepository(
			d.PostgresDatabase(ctx),
		)
	}

	return d.ordersRepository
}

func (d *diContainer) OrdersService(ctx context.Context) api.OrdersService {
	if d.ordersService == nil {
		d.ordersService = orders.NewService(
			d.OrdersRepository(ctx),
			d.PaymentClient(),
			d.InventoryClient(),
			d.ProducerService(),
		)
	}

	return d.ordersService
}

func (d *diContainer) ServerImplementation(ctx context.Context) *api.OrderImplementation {
	if d.serverImplementation == nil {
		d.serverImplementation = api.NewOrderHandler(
			d.OrdersService(ctx),
			d.PaymentClient(),
			d.InventoryClient(),
		)
	}

	return d.serverImplementation
}
