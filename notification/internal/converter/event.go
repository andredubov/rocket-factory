package converter

import (
	"github.com/andredubov/rocket-factory/notification/internal/model"
	events_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/events/v1"
)

func OrderAssembledEventToProtobufEvent(event model.OrderAssembledEvent) *events_v1.ShipAssembled {
	return &events_v1.ShipAssembled{
		EventUuid:    event.UUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}
}

func OrderPaidEventToProtobufEvent(event model.OrderPaidEvent) *events_v1.OrderPaid {
	return &events_v1.OrderPaid{
		EventUuid:       event.UUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   string(event.PaymentMethod),
		TransactionUuid: event.TransactionUUID,
	}
}

func OrderPaidEventFromProtobufEvent(event *events_v1.OrderPaid) (model.OrderPaidEvent, error) {
	paymentMethod, err := model.NewPaymentMethod(event.PaymentMethod)
	if err != nil {
		return model.OrderPaidEvent{}, model.ErrInvalidPaymentMethod
	}

	return model.OrderPaidEvent{
		UUID:            event.EventUuid,
		OrderUUID:       event.OrderUuid,
		UserUUID:        event.UserUuid,
		PaymentMethod:   paymentMethod,
		TransactionUUID: event.TransactionUuid,
	}, nil
}

func OrderAssembledEventFromProtobufEvent(event *events_v1.ShipAssembled) model.OrderAssembledEvent {
	return model.OrderAssembledEvent{
		UUID:         event.EventUuid,
		OrderUUID:    event.OrderUuid,
		UserUUID:     event.UserUuid,
		BuildTimeSec: event.BuildTimeSec,
	}
}
