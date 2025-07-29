package tests

import (
	"testing"

	"github.com/dvln/testify/suite"

	api "github.com/andredubov/rocket-factory/inventory/internal/api/v1/inventory"
	"github.com/andredubov/rocket-factory/inventory/internal/service/inventory"
	"github.com/andredubov/rocket-factory/inventory/internal/service/mocks"
)

// InventoryServiceSuite defines the test suite for inventory service integration tests.
// Contains mocked dependencies and the service instance being tested.
type InventoryServiceSuite struct {
	suite.Suite
	inventoryRepository *mocks.InventoryRepository
	inventoryService    api.InventoryService
}

// SetupTest initializes the test environment before each test case.
func (s *InventoryServiceSuite) SetupTest() {
	s.inventoryRepository = mocks.NewInventoryRepository(s.T())
	s.inventoryService = inventory.NewService(s.inventoryRepository)
}

// TearDownTest performs cleanup after each test case.
func (s *InventoryServiceSuite) TearDownTest() {
}

// TestInventoryServiceIntegration is the entry point for running the inventory service test suite.
func TestInventoryServiceIntegration(t *testing.T) {
	suite.Run(t, new(InventoryServiceSuite))
}
