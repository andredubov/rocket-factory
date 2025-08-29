package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/andredubov/rocket-factory/assembly/internal/model"
	producer "github.com/andredubov/rocket-factory/assembly/internal/service/producer/order_assembled_producer"
	"github.com/andredubov/rocket-factory/platform/pkg/kafka/mocks"
	events_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/events/v1"
)

func TestNewService(t *testing.T) {
	// Setup
	mockProducer := mocks.NewProducer(t)

	// Test
	service := producer.NewService(mockProducer)

	// Verify
	assert.NotNil(t, service)
}

func TestProduceOrderAssembledEvent_Success(t *testing.T) {
	// Setup
	var (
		mockProducer = mocks.NewProducer(t)
		service      = producer.NewService(mockProducer)
		ctx          = context.Background()
		event        = model.OrderAssembledEvent{
			UUID: gofakeit.UUID(),
		}
		expectedMsg = &events_v1.ShipAssembled{
			EventUuid: event.UUID,
		}
	)

	expectedPayload, err := proto.Marshal(expectedMsg)
	assert.NoError(t, err)

	// Expectation
	mockProducer.On("Send", ctx, []byte(event.UUID), expectedPayload).Return(nil)

	// Test
	err = service.ProduceOrderAssembledEvent(ctx, event)

	// Verify
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestProduceOrderAssembledEvent_WithMockExpecter(t *testing.T) {
	// Setup
	var (
		mockProducer = mocks.NewProducer(t)
		service      = producer.NewService(mockProducer)
		ctx          = context.Background()
		event        = model.OrderAssembledEvent{
			UUID: gofakeit.UUID(),
		}
		expectedMsg = &events_v1.ShipAssembled{
			EventUuid: event.UUID,
		}
	)

	expectedPayload, err := proto.Marshal(expectedMsg)
	assert.NoError(t, err)

	// Expectation
	mockProducer.EXPECT().Send(ctx, []byte(event.UUID), expectedPayload).Return(nil)

	// Test
	err = service.ProduceOrderAssembledEvent(ctx, event)

	// Verify
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestService_ProduceOrderAssembledEvent_Success(t *testing.T) {
	// Setup
	var (
		mockProducer = mocks.NewProducer(t)
		svc          = producer.NewService(mockProducer)
		ctx          = context.Background()
		event        = model.OrderAssembledEvent{
			UUID: gofakeit.UUID(),
		}
		expectedMsg = &events_v1.ShipAssembled{
			EventUuid: event.UUID,
		}
	)

	expectedPayload, err := proto.Marshal(expectedMsg)
	assert.NoError(t, err)

	// Expectation
	mockProducer.On("Send", ctx, []byte(event.UUID), expectedPayload).Return(nil)

	// Test
	err = svc.ProduceOrderAssembledEvent(ctx, event)

	// Verify
	assert.NoError(t, err)
	mockProducer.AssertExpectations(t)
}
