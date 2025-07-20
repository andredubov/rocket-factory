package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/andredubov/rocket-factory/order/internal/repository/mocks"
	"github.com/andredubov/rocket-factory/order/internal/service"
	orders "github.com/andredubov/rocket-factory/order/internal/service/order"
)

// OrdersServiceSuite defines the test suite for order service integration tests.
type OrdersServiceSuite struct {
	suite.Suite
	ordersRepository *mocks.Orders
	ordersService    service.Orders
}

// SetupTest initializes the test environment before each test case.
func (s *OrdersServiceSuite) SetupTest() {
	s.ordersRepository = mocks.NewOrders(s.T())
	s.ordersService = orders.NewService(s.ordersRepository)
}

// TearDownTest performs cleanup after each test case.
func (s *OrdersServiceSuite) TearDownTest() {
}

// TestOrdersServiceIntegration is the entry point for running the inventory service test suite.
func TestOrdersServiceIntegration(t *testing.T) {
	suite.Run(t, new(OrdersServiceSuite))
}
