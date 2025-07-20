package tests

import (
	"testing"

	"github.com/dvln/testify/suite"
)

// OrdersRepoConverterSuite defines the test suite for order repository converter integration tests.
type OrdersRepoConverterSuite struct {
	suite.Suite
}

// SetupTest prepares the test environment before each test case execution.
func (s *OrdersRepoConverterSuite) SetupTest() {
}

// TearDownTest performs cleanup after each test case execution.
func (s *OrdersRepoConverterSuite) TearDownTest() {
}

// TestOrdersRepoConverterIntegration is the entry point for running the converter test suite.
func TestOrdersRepoConveterIntegration(t *testing.T) {
	suite.Run(t, new(OrdersRepoConverterSuite))
}
