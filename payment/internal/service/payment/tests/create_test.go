package tests

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/andredubov/rocket-factory/payment/internal/model"
)

func (s *PaymentServiceSuite) TestCreateSuccess() {
	var (
		ctx = context.Background()

		payment = model.Payment{
			UserUuid:      gofakeit.UUID(),
			OrderUuid:     gofakeit.UUID(),
			PaymentMethod: model.PaymentMethodCard,
		}
	)

	uuid, err := s.paymentService.Create(ctx, payment)

	s.Require().NoError(err)
	s.Require().Nil(err)
	s.Require().NotEmpty(uuid)
}
