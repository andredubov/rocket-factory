package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/andredubov/rocket-factory/notification/internal/converter"
	"github.com/andredubov/rocket-factory/notification/internal/model"
	events_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/events/v1"
)

type orderPaidEventDecoder struct{}

func NewOrderPaidEventDecoder() *orderPaidEventDecoder {
	return &orderPaidEventDecoder{}
}

func (d *orderPaidEventDecoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	pb := events_v1.OrderPaid{}
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	event, err := converter.OrderPaidEventFromProtobufEvent(&pb)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to convert OrderPaidEvent from protobuf event: %w", err)
	}

	return event, nil
}
