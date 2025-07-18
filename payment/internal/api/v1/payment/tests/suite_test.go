package test

import (
	"testing"

	"github.com/dvln/testify/suite"

	server "github.com/andredubov/rocket-factory/payment/internal/api/v1/payment"
	"github.com/andredubov/rocket-factory/payment/internal/service/mocks"
)

type APISuite struct {
	suite.Suite
	paymentService *mocks.Payments
	grpcServer     *server.PaymentImplementation
}

func (s *APISuite) SetupTest() {
	s.paymentService = mocks.NewPayments(s.T())
	s.grpcServer = server.NewPaymentImplementation(s.paymentService)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
