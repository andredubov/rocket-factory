package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/payment/internal/model"
)

// Create implements the Payments interface by generating a new payment transaction.
func (p *paymentService) Create(ctx context.Context, payment model.Payment) (string, error) {
	uuid := uuid.New().String()
	return uuid, nil
}
