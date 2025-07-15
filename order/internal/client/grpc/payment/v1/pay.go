package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/andredubov/rocket-factory/order/internal/client/converter"
	"github.com/andredubov/rocket-factory/order/internal/model"
)

func (c *paymentClient) PayOrder(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	// Создание запроса в платежный сервис
	paymentRequest := converter.OrderToPayOrderRequest(order)

	// Вызов платежного сервиса
	paymentResponse, err := c.generatedClient.PayOrder(ctx, paymentRequest)
	if err != nil {
		return uuid.New(), fmt.Errorf("payment service error: %w", err)
	}

	// Извлечение UUID транзакции из ответа платежного сервиса
	transactionUUID, err := converter.TransactionUuidFromResponse(paymentResponse)
	if err != nil {
		return uuid.New(), fmt.Errorf("invalid transaction uuid: %w", err)
	}

	return transactionUUID, nil
}
