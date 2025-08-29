package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderAssembledEventConsumerEnvConfig struct {
	Topic   string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type orderAssembledEventConsumerConfig struct {
	raw orderAssembledEventConsumerEnvConfig
}

func NewOrderAssembledEventConsumerConfig() (*orderAssembledEventConsumerConfig, error) {
	var raw orderAssembledEventConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderAssembledEventConsumerConfig{raw: raw}, nil
}

func (cfg *orderAssembledEventConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *orderAssembledEventConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *orderAssembledEventConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
