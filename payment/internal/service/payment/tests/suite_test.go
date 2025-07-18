package tests

import (
	"testing"

	"github.com/dvln/testify/suite"

	"github.com/andredubov/rocket-factory/payment/internal/service"
	"github.com/andredubov/rocket-factory/payment/internal/service/payment"
)

type PaymentServiceSuite struct {
	suite.Suite
	paymentService service.Payments
}

func (s *PaymentServiceSuite) SetupTest() {
	s.paymentService = payment.NewService()
}

func (s *PaymentServiceSuite) TearDownTest() {
}

func TestPaymentServiceIntegration(t *testing.T) {
	suite.Run(t, new(PaymentServiceSuite))
}
