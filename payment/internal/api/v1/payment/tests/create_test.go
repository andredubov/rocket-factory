package test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	api "github.com/andredubov/rocket-factory/payment/internal/api/v1/payment"
	"github.com/andredubov/rocket-factory/payment/internal/api/v1/payment/mocks"
	"github.com/andredubov/rocket-factory/payment/internal/converter"
	"github.com/andredubov/rocket-factory/payment/internal/model"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

func TestCreateSuccess(t *testing.T) {
	var (
		paymentService = mocks.NewPaymentService(t)
		grpcServer     = api.NewPaymentImplementation(paymentService)
		ctx            = context.Background()
		expectedUuid   = gofakeit.UUID()

		payment = model.Payment{
			OrderUuid:     gofakeit.UUID(),
			UserUuid:      gofakeit.UUID(),
			PaymentMethod: model.PaymentMethodCard,
		}

		req = &payment_v1.PayOrderRequest{
			OrderUuid:     payment.OrderUuid,
			UserUuid:      payment.UserUuid,
			PaymentMethod: payment_v1.PaymentMethod(payment.PaymentMethod),
		}
	)

	paymentService.On("Create", ctx, converter.PaymentFromRequest(req)).Return(expectedUuid, nil)

	res, err := grpcServer.PayOrder(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expectedUuid, res.TransactionUuid)
}
