package converter

import (
	"fmt"

	"github.com/google/uuid"

	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

func TransactionUuidFromResponse(response *payment_v1.PayOrderResponse) (uuid.UUID, error) {
	// Парсинг UUID транзакции
	transactionUUID, err := uuid.Parse(response.GetTransactionUuid())
	if err != nil {
		return uuid.New(), fmt.Errorf("invalid transaction uuid: %w", err)
	}

	return transactionUUID, nil
}
