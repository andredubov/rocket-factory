package tests

import (
	"testing"

	"github.com/dvln/testify/suite"
)

type InventoryRepoConverterSuite struct {
	suite.Suite
}

func (s *InventoryRepoConverterSuite) SetupTest() {
}

func (s *InventoryRepoConverterSuite) TearDownTest() {
}

func TestInventoryRepoConveterIntegration(t *testing.T) {
	suite.Run(t, new(InventoryRepoConverterSuite))
}
