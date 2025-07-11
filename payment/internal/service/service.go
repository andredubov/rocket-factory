package service

import (
	"context"

	"github.com/andredubov/rocket-factory/payment/internal/service/model"
)

// Payments defines the contract for payment processing operations.
// Implementations of this interface are responsible for handling the complete
// payment lifecycle in the system.
type Payments interface {
	Create(ctx context.Context, payment model.Payment) (string, error)
}
