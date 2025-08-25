package converter

import (
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/andredubov/rocket-factory/order/internal/model"
	events_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/events/v1"
)

func OrderPaidEventToProtobufEvent(event model.OrderPaidEvent) *events_v1.OrderPaid {
	return &events_v1.OrderPaid{
		EventUuid:       event.UUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   string(event.PaymentMethod),
		TransactionUuid: event.TrasactionUUID,
	}
}

func OrderAssembledEventFromProtobufEvent(event *events_v1.ShipAssembled) model.OrderAssembledEvent {
	return model.OrderAssembledEvent{
		UUID:         event.EventUuid,
		OrderUUID:    event.OrderUuid,
		UserUUID:     event.UserUuid,
		BuildTimeSec: event.BuildTimeSec,
	}
}

func OrderToOrderPaidEvent(order *model.Order) model.OrderPaidEvent {
	return model.OrderPaidEvent{
		UUID:           order.OrderUUID.String(),
		UserUUID:       order.UserUUID.String(),
		OrderUUID:      order.OrderUUID.String(),
		PaymentMethod:  order.PaymentInfo.PaymentMethod,
		TrasactionUUID: order.PaymentInfo.TransactionUUID.String(),
	}
}

func OrderAssembledEventToOrderUpdateInfo(event model.OrderAssembledEvent) (model.OrderUpdateInfo, error) {
	orderUUID, err := uuid.Parse(event.OrderUUID)
	if err != nil {
		return model.OrderUpdateInfo{}, err
	}

	updateInfo := model.OrderUpdateInfo{
		OrderUUID: orderUUID,
		Status:    lo.ToPtr(model.OrderStatusAssembled),
	}

	return updateInfo, nil
}
