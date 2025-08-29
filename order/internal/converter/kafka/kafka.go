package kafka

import (
	"github.com/andredubov/rocket-factory/order/internal/model"
)

type OrderAssembledEventDecoder interface {
	Decode(data []byte) (model.OrderAssembledEvent, error)
}
