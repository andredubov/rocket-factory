package service

import (
	"context"

	"github.com/andredubov/rocket-factory/assembly/internal/model"
)

// ConsumerService defines the interface for service that consunes order paid events
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

// ProducerService defines the interface for service that produces order assembled events
type ProducerService interface {
	ProduceOrderAssembledEvent(ctx context.Context, event model.OrderAssembledEvent) error
}
