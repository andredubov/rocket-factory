package tests

import (
	"testing"

	"github.com/dvln/testify/suite"

	"github.com/andredubov/rocket-factory/order/internal/repository/order/memory"
	"github.com/andredubov/rocket-factory/order/internal/service"
)

// InventoryRepositorySuite is a test suite for orders repository operations.
type OrdersRepositorySuite struct {
	suite.Suite
	ordersRepository service.OrdersRepository
}

// SetupTest prepares the test environment before each test case execution.
func (s *OrdersRepositorySuite) SetupTest() {
	s.ordersRepository = memory.NewOrderRepository()
}

// TearDownTest performs cleanup after each test case.
func (s *OrdersRepositorySuite) TearDownTest() {
}

// TestInventoryRepositoryIntegration is the entry point for running the repository test suite.
func TestInventoryRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(OrdersRepositorySuite))
}
