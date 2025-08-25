package decoder

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/andredubov/rocket-factory/assembly/internal/converter"
	"github.com/andredubov/rocket-factory/assembly/internal/model"
	events_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderPaidEventDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	pb := events_v1.OrderPaid{}
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return converter.OrderPaidEventFromProtobufEvent(&pb), nil
}
