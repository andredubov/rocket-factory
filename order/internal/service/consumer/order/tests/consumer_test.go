package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	kafkaConverterMocks "github.com/andredubov/rocket-factory/order/internal/converter/kafka/mocks"
	"github.com/andredubov/rocket-factory/order/internal/service/consumer/order"
	repoMocks "github.com/andredubov/rocket-factory/order/internal/service/mocks"
	"github.com/andredubov/rocket-factory/platform/pkg/kafka/mocks"
)

func TestNewService(t *testing.T) {
	// Создаем моки для зависимостей
	mockConsumer := mocks.NewConsumer(t)
	mockDecoder := kafkaConverterMocks.NewOrderAssembledEventDecoder(t)
	mockRepo := repoMocks.NewOrdersRepository(t)

	// Вызываем тестируемую функцию
	service := order.NewService(mockConsumer, mockDecoder, mockRepo)

	// Проверяем, что сервис создан корректно
	assert.NotNil(t, service)
}
