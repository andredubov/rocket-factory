package env

func NewOrderPaidEventConsumerConfig() (*KafkaConsumerConfig, error) {
	return newKafkaConsumerConfig(
		"ORDER_PAID_TOPIC_NAME",
		"ORDER_PAID_CONSUMER_GROUP_ID",
	)
}
