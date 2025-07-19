package tests

import (
	"testing"

	"github.com/dvln/testify/suite"

	"github.com/andredubov/rocket-factory/inventory/internal/repository/mocks"
	"github.com/andredubov/rocket-factory/inventory/internal/service"
	"github.com/andredubov/rocket-factory/inventory/internal/service/inventory"
)

// InventoryServiceSuite defines the test suite for inventory service integration tests.
// Contains mocked dependencies and the service instance being tested.
type InventoryServiceSuite struct {
	suite.Suite
	inventoryRepository *mocks.Inventory
	inventoryService    service.Inventory
}

// SetupTest initializes the test environment before each test case.
func (s *InventoryServiceSuite) SetupTest() {
	s.inventoryRepository = mocks.NewInventory(s.T())
	s.inventoryService = inventory.NewService(s.inventoryRepository)
}

// TearDownTest performs cleanup after each test case.
func (s *InventoryServiceSuite) TearDownTest() {
}

// TestInventoryServiceIntegration is the entry point for running the inventory service test suite.
func TestInventoryServiceIntegration(t *testing.T) {
	suite.Run(t, new(InventoryServiceSuite))
}
