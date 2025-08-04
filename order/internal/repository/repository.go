package repository

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

func ErrOrderAlreadyExistsWith(uuid uuid.UUID) error {
	return fmt.Errorf("%w: %s", model.ErrOrderAlreadyExists, uuid)
}

func ErrOrderNotFoundWith(uuid uuid.UUID) error {
	return fmt.Errorf("%w: %s", model.ErrOrderNotFound, uuid)
}
