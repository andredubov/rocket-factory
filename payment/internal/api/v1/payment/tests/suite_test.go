package test

import (
	"testing"

	"github.com/dvln/testify/suite"

	api "github.com/andredubov/rocket-factory/payment/internal/api/v1/payment"
	"github.com/andredubov/rocket-factory/payment/internal/api/v1/payment/mocks"
)

type APISuite struct {
	suite.Suite
	paymentService *mocks.PaymentService
	grpcServer     *api.PaymentImplementation
}

func (s *APISuite) SetupTest() {
	s.paymentService = mocks.NewPaymentService(s.T())
	s.grpcServer = api.NewPaymentImplementation(s.paymentService)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
