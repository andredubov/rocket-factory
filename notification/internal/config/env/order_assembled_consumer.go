package env

func NewOrderAssembledEventConsumerConfig() (*KafkaConsumerConfig, error) {
	return newKafkaConsumerConfig(
		"ORDER_ASSEMBLED_TOPIC_NAME",
		"ORDER_ASSEMBLED_CONSUMER_GROUP_ID",
	)
}
