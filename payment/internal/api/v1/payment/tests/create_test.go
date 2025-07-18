package test

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/andredubov/rocket-factory/payment/internal/converter"
	"github.com/andredubov/rocket-factory/payment/internal/model"
	payment_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/payment/v1"
)

func (s *APISuite) TestCreateSuccess() {
	var (
		ctx          = context.Background()
		expectedUuid = gofakeit.UUID()

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

	s.paymentService.On("Create", ctx, converter.PaymentFromRequest(req)).Return(expectedUuid, nil)

	res, err := s.grpcServer.PayOrder(ctx, req)
	s.Require().NoError(err)
	s.Require().Nil(err)
	s.Require().Equal(expectedUuid, res.TransactionUuid)
}
