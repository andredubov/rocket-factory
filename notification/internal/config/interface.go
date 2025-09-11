package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OTLPEndpoint() string
	ServiceName() string
	ServiceEnvironment() string
}

type KafkaConfig interface {
	Brokers() []string
}

type TelegramConfig interface {
	TelegramBotToken() string
	TelegramChatID() int64
}

type OrderEventConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
