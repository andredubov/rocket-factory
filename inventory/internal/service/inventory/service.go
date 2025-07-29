package inventory

import (
	api "github.com/andredubov/rocket-factory/inventory/internal/api/v1/inventory"
	"github.com/andredubov/rocket-factory/inventory/internal/service"
)

type inventoryService struct {
	inventoryRepository service.InventoryRepository
}

func NewService(repo service.InventoryRepository) api.InventoryService {
	return &inventoryService{
		inventoryRepository: repo,
	}
}
