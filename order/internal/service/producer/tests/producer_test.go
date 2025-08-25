package tests

import (
	"context"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/andredubov/rocket-factory/order/internal/converter"
	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/service/producer"
	"github.com/andredubov/rocket-factory/platform/pkg/kafka/mocks"
)

func TestNewService(t *testing.T) {
	mockProducer := mocks.NewProducer(t)
	service := producer.NewService(mockProducer)

	assert.NotNil(t, service)
}

func TestProduceOrderPaidEvent_Success(t *testing.T) {
	mockProducer := mocks.NewProducer(t)
	service := producer.NewService(mockProducer)

	ctx := context.Background()
	event := model.OrderPaidEvent{
		UUID: "test-uuid-123",
	}

	// Конвертируем событие в protobuf для получения ожидаемого payload
	expectedMsg := converter.OrderPaidEventToProtobufEvent(event)
	expectedPayload, err := proto.Marshal(expectedMsg)
	assert.NoError(t, err)

	// Настраиваем ожидание вызова Send
	mockProducer.On("Send", ctx, []byte(event.UUID), expectedPayload).
		Return(nil)

	// Вызываем тестируемый метод
	err = service.ProduceOrderPaidEvent(ctx, event)

	// Проверяем результаты
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestProduceOrderPaidEvent_WithMockExpecter(t *testing.T) {
	mockProducer := mocks.NewProducer(t)
	service := producer.NewService(mockProducer)

	ctx := context.Background()
	event := model.OrderPaidEvent{
		UUID: "test-uuid-456",
	}

	// Конвертируем событие в protobuf
	expectedMsg := converter.OrderPaidEventToProtobufEvent(event)
	expectedPayload, err := proto.Marshal(expectedMsg)
	assert.NoError(t, err)

	// Используем Expecter для более типобезопасного подхода
	mockProducer.EXPECT().
		Send(ctx, []byte(event.UUID), expectedPayload).
		Return(nil)

	// Вызываем тестируемый метод
	err = service.ProduceOrderPaidEvent(ctx, event)

	// Проверяем результаты
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestService_ProduceOrderPaidEvent_Success(t *testing.T) {
	mockProducer := mocks.NewProducer(t)
	svc := producer.NewService(mockProducer)

	ctx := context.Background()
	event := model.OrderPaidEvent{
		UUID: "order-uuid-1",
	}

	expectedMsg := converter.OrderPaidEventToProtobufEvent(event)
	expectedPayload, err := proto.Marshal(expectedMsg)
	assert.NoError(t, err)

	mockProducer.On("Send", ctx, []byte(event.UUID), expectedPayload).Return(nil)

	err = svc.ProduceOrderPaidEvent(ctx, event)
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}
