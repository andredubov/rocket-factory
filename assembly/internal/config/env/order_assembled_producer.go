package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderAssembledEventProducerEnvConfig struct {
	TopicName string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
}

type orderAssembledEventProducerConfig struct {
	raw orderAssembledEventProducerEnvConfig
}

func NewOrderAssembledEventProducerConfig() (*orderAssembledEventProducerConfig, error) {
	var raw orderAssembledEventProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderAssembledEventProducerConfig{raw: raw}, nil
}

func (cfg *orderAssembledEventProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

// Config возвращает конфигурацию для sarama consumer
func (cfg *orderAssembledEventProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
