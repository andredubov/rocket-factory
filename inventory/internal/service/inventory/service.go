package inventory

import (
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/service"
)

type inventoryService struct {
	inventoryRepository repository.Inventory
}

func NewService(repo repository.Inventory) service.Inventory {
	return &inventoryService{
		inventoryRepository: repo,
	}
}
