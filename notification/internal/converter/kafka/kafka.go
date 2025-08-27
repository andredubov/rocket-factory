package converter

import "github.com/andredubov/rocket-factory/notification/internal/model"

type OrderPaidEventDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}

type OrderAssembledEventDecoder interface {
	Decode(data []byte) (model.OrderAssembledEvent, error)
}
