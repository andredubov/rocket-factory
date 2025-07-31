package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"

	api "github.com/andredubov/rocket-factory/order/internal/api/v1/order"
	"github.com/andredubov/rocket-factory/order/internal/service/mocks"
	orders "github.com/andredubov/rocket-factory/order/internal/service/order"
)

// OrdersServiceSuite defines the test suite for order service integration tests.
type OrdersServiceSuite struct {
	suite.Suite
	ordersRepository *mocks.OrdersRepository
	paymentClient    *mocks.PaymentClient
	inventoryClient  *mocks.InventoryClient
	ordersService    api.OrdersService
}

// SetupTest initializes the test environment before each test case.
func (s *OrdersServiceSuite) SetupTest() {
	s.ordersRepository = mocks.NewOrdersRepository(s.T())
	s.paymentClient = mocks.NewPaymentClient(s.T())
	s.inventoryClient = mocks.NewInventoryClient(s.T())
	s.ordersService = orders.NewService(s.ordersRepository, s.paymentClient, s.inventoryClient)
}

// TearDownTest performs cleanup after each test case.
func (s *OrdersServiceSuite) TearDownTest() {
}

// TestOrdersServiceIntegration is the entry point for running the inventory service test suite.
func TestOrdersServiceIntegration(t *testing.T) {
	suite.Run(t, new(OrdersServiceSuite))
}
