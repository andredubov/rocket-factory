package tests

import (
	"testing"

	"github.com/dvln/testify/suite"

	api "github.com/andredubov/rocket-factory/payment/internal/api/v1/payment"
	"github.com/andredubov/rocket-factory/payment/internal/service/payment"
)

type PaymentServiceSuite struct {
	suite.Suite
	paymentService api.PaymentService
}

func (s *PaymentServiceSuite) SetupTest() {
	s.paymentService = payment.NewService()
}

func (s *PaymentServiceSuite) TearDownTest() {
}

func TestPaymentServiceIntegration(t *testing.T) {
	suite.Run(t, new(PaymentServiceSuite))
}
