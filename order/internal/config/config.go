package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/andredubov/rocket-factory/order/internal/config/env"
	iamClient "github.com/andredubov/rocket-factory/order/internal/config/env/iam"
	inventoryClient "github.com/andredubov/rocket-factory/order/internal/config/env/inventory"
	paymentClient "github.com/andredubov/rocket-factory/order/internal/config/env/payment"
)

var appConfig *config

type config struct {
	Logger                      LoggerConfig
	HTTPServer                  HTTPConfig
	PostgresDB                  PostgresDBConfig
	InventoryClient             GRPCConfig
	PaymentClient               GRPCConfig
	IAMClient                   GRPCConfig
	Kafka                       KafkaConfig
	OrderPaidEventProducer      OrderPaidEventProducerConfig
	OrderAssembledEventConsumer OrderAssembledEventConsumerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	httpConfig, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}

	postgresDBConfig, err := env.NewPostgresDBConfig()
	if err != nil {
		return err
	}

	invetoryClientConfig, err := inventoryClient.NewGRPCConfig()
	if err != nil {
		return err
	}

	paymentClientConfig, err := paymentClient.NewGRPCConfig()
	if err != nil {
		return err
	}

	iamClientConfig, err := iamClient.NewGRPCConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidEventProducerCfg, err := env.NewOrderPaidEventProducerConfig()
	if err != nil {
		return err
	}

	orderAssembledEventConsumerCfg, err := env.NewOrderAssembledEventConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                      loggerConfig,
		HTTPServer:                  httpConfig,
		PostgresDB:                  postgresDBConfig,
		InventoryClient:             invetoryClientConfig,
		PaymentClient:               paymentClientConfig,
		IAMClient:                   iamClientConfig,
		Kafka:                       kafkaCfg,
		OrderPaidEventProducer:      orderPaidEventProducerCfg,
		OrderAssembledEventConsumer: orderAssembledEventConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
