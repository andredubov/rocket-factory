package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GRPCConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
	ReadHeaderTimeout() time.Duration
}

type PostgresDBConfig interface {
	DSN() string
	MigrationDirectory() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidEventProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderAssembledEventConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
