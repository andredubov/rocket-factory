package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/andredubov/rocket-factory/notification/internal/converter"
	"github.com/andredubov/rocket-factory/notification/internal/model"
	events_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/events/v1"
)

type orderAssembledEventDecoder struct{}

func NewOrderAssembledEventDecoder() *orderAssembledEventDecoder {
	return &orderAssembledEventDecoder{}
}

func (d *orderAssembledEventDecoder) Decode(data []byte) (model.OrderAssembledEvent, error) {
	pb := events_v1.ShipAssembled{}
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return converter.OrderAssembledEventFromProtobufEvent(&pb), nil
}
