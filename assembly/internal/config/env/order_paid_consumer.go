package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderPaidEventConsumerEnvConfig struct {
	Topic   string `env:"ORDER_PAID_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

type orderPaidEventConsumerConfig struct {
	raw orderPaidEventConsumerEnvConfig
}

func NewOrderPaidEventConsumerConfig() (*orderPaidEventConsumerConfig, error) {
	var raw orderPaidEventConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidEventConsumerConfig{raw: raw}, nil
}

func (cfg *orderPaidEventConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *orderPaidEventConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *orderPaidEventConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
