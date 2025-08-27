package service

import (
	"context"

	"github.com/andredubov/rocket-factory/notification/internal/model"
)

// ConsumerService defines the interface for service that consunes order paid events
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

// TelegramService defines the interface for service that send notification to Telegram
type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, uuid string, event model.OrderPaidEvent) error
	SendOrderAssembledNotification(ctx context.Context, uuid string, event model.OrderAssembledEvent) error
}
