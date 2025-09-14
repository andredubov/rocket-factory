package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/payment/internal/model"
	"github.com/andredubov/rocket-factory/payment/internal/service/payment"
)

func TestCreateSuccess(t *testing.T) {
	var (
		paymentService = payment.NewService()
		ctx            = context.Background()
		payment        = model.Payment{
			UserUuid:      gofakeit.UUID(),
			OrderUuid:     gofakeit.UUID(),
			PaymentMethod: model.PaymentMethodCard,
		}
	)

	uuid, err := paymentService.Create(ctx, payment)

	require.NoError(t, err)
	require.Nil(t, err)
	require.NotEmpty(t, uuid)
}
