package kafka

import (
	"github.com/andredubov/rocket-factory/assembly/internal/model"
)

type OrderPaidEventDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
