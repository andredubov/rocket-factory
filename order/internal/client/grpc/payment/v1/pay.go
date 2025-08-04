package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/client/converter"
	"github.com/andredubov/rocket-factory/order/internal/model"
)

// PayOrder processes payment for the given order through the payment service.
func (c *paymentClient) PayOrder(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	paymentRequest := converter.OrderToPayOrderRequest(order)

	paymentResponse, err := c.generatedClient.PayOrder(ctx, paymentRequest)
	if err != nil {
		return uuid.New(), fmt.Errorf("payment service error: %w", err)
	}

	transactionUUID, err := converter.TransactionUuidFromResponse(paymentResponse)
	if err != nil {
		return uuid.New(), fmt.Errorf("invalid transaction uuid: %w", err)
	}

	return transactionUUID, nil
}
