package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderPaidEventProducerEnvConfig struct {
	TopicName string `env:"ORDER_PAID_TOPIC_NAME,required"`
}

type orderPaidEventProducerConfig struct {
	raw orderPaidEventProducerEnvConfig
}

func NewOrderPaidEventProducerConfig() (*orderPaidEventProducerConfig, error) {
	var raw orderPaidEventProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidEventProducerConfig{raw: raw}, nil
}

func (cfg *orderPaidEventProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

// Config возвращает конфигурацию для sarama consumer
func (cfg *orderPaidEventProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
