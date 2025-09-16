package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPEndpoint() string
	ServiceName() string
	ServiceEnvironment() string
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

type MetricsConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
	CollectorTimeout() time.Duration
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}
