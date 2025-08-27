package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/andredubov/rocket-factory/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                      LoggerConfig
	Kafka                       KafkaConfig
	OrderPaidEventConsumer      OrderEventConsumerConfig
	OrderAssembledEventConsumer OrderEventConsumerConfig
	Telegram                    TelegramConfig
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

	orderPaidEventConsumerCfg, err := env.NewOrderPaidEventConsumerConfig()
	if err != nil {
		return err
	}

	orderAssembledEventConsumerCfg, err := env.NewOrderAssembledEventConsumerConfig()
	if err != nil {
		return err
	}

	telegramCfg, err := env.NewTelegramConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                      loggerConfig,
		Kafka:                       kafkaCfg,
		OrderPaidEventConsumer:      orderPaidEventConsumerCfg,
		OrderAssembledEventConsumer: orderAssembledEventConsumerCfg,
		Telegram:                    telegramCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
