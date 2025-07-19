package tests

import (
	"testing"

	"github.com/dvln/testify/suite"
)

type InventoryRepoModelsSuite struct {
	suite.Suite
}

func (s *InventoryRepoModelsSuite) SetupTest() {
}

func (s *InventoryRepoModelsSuite) TearDownTest() {
}

func TestInventoryRepoModelsIntegration(t *testing.T) {
	suite.Run(t, new(InventoryRepoModelsSuite))
}
