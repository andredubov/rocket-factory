package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dvln/testify/suite"

	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/memory"
)

// InventoryRepositorySuite is a test suite for inventory repository operations.
type InventoryRepositorySuite struct {
	suite.Suite
	inventoryRepository repository.Inventory
	parts               []model.Part
}

// SetupTest prepares the test environment before each test case execution.
func (s *InventoryRepositorySuite) SetupTest() {
	s.inventoryRepository = memory.NewInventoryRepository()
	s.parts = make([]model.Part, 0)

	// Seed the repository with test data
	var (
		ctx        = context.Background()
		countries  = []string{"USA", "Germany", "Japan", "China"}
		categories = []model.PartCategory{
			model.PartCategoryEngine,
			model.PartCategoryFuel,
			model.PartCategoryPorthole,
			model.PartCategoryWing,
		}
	)

	for i := 0; i < 20; i++ {
		part := model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          "Part-" + gofakeit.Word(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Float64Range(1, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      categories[gofakeit.IntRange(0, 3)],
			Dimensions: model.Dimensions{
				Length: gofakeit.Float64Range(1, 100),
				Width:  gofakeit.Float64Range(1, 100),
				Height: gofakeit.Float64Range(1, 100),
				Weight: gofakeit.Float64Range(1, 100),
			},
			Manufacturer: model.Manufacturer{
				Name:    gofakeit.Company(),
				Country: countries[gofakeit.IntRange(0, 3)],
				Website: gofakeit.URL(),
			},
			Tags:      []string{"tag-" + gofakeit.Word(), "tag-" + gofakeit.Word()},
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}
		err := s.inventoryRepository.AddPart(ctx, part)
		s.Require().NoError(err)
		s.parts = append(s.parts, part)
	}
}

// TearDownTest performs cleanup after each test case.
func (s *InventoryRepositorySuite) TearDownTest() {
}

// TestInventoryRepositoryIntegration is the entry point for running the repository test suite.
func TestInventoryRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(InventoryRepositorySuite))
}
