package env

import (
	"os"

	"github.com/IBM/sarama"
)

type KafkaConsumerConfig struct {
	topic   string
	groupID string
}

func newKafkaConsumerConfig(topicEnvVar, groupIDEnvVar string) (*KafkaConsumerConfig, error) {
	topic, topicOk := os.LookupEnv(topicEnvVar)
	groupID, groupIDOk := os.LookupEnv(groupIDEnvVar)

	if !topicOk || topic == "" {
		return nil, &envError{msg: "environment variable " + topicEnvVar + " is required"}
	}
	if !groupIDOk || groupID == "" {
		return nil, &envError{msg: "environment variable " + groupIDEnvVar + " is required"}
	}

	return &KafkaConsumerConfig{
		topic:   topic,
		groupID: groupID,
	}, nil
}

func (cfg *KafkaConsumerConfig) Topic() string {
	return cfg.topic
}

func (cfg *KafkaConsumerConfig) GroupID() string {
	return cfg.groupID
}

func (cfg *KafkaConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}

type envError struct {
	msg string
}

func (e *envError) Error() string {
	return e.msg
}
