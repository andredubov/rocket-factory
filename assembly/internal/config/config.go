package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/andredubov/rocket-factory/assembly/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                      LoggerConfig
	Kafka                       KafkaConfig
	OrderAssembledEventProducer OrderAssembledProducerConfig
	OrderPaidEventConsumer      OrderPaidEventConsumerConfig
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

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderAssembledEventProducerCfg, err := env.NewOrderAssembledEventProducerConfig()
	if err != nil {
		return err
	}

	orderPaidEventConsumerCfg, err := env.NewOrderPaidEventConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                      loggerConfig,
		Kafka:                       kafkaCfg,
		OrderAssembledEventProducer: orderAssembledEventProducerCfg,
		OrderPaidEventConsumer:      orderPaidEventConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
