package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	kafkaConverterMocks "github.com/andredubov/rocket-factory/assembly/internal/converter/kafka/mocks"
	"github.com/andredubov/rocket-factory/assembly/internal/service/consumer/order"
	"github.com/andredubov/rocket-factory/assembly/internal/service/producer"
	"github.com/andredubov/rocket-factory/platform/pkg/kafka/mocks"
)

func TestNewService(t *testing.T) {
	// Setup
	var (
		mockConsumer    = mocks.NewConsumer(t)
		mockDecoder     = kafkaConverterMocks.NewOrderPaidEventDecoder(t)
		mockProducer    = mocks.NewProducer(t)
		producerService = producer.NewService(mockProducer)
	)

	// Test
	service := order.NewService(mockConsumer, mockDecoder, producerService)

	// Verify
	assert.NotNil(t, service)
}
