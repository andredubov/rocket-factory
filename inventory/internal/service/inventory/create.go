package inventory

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
)

func (i *inventoryService) AddPart(ctx context.Context, part model.Part) error {
	return i.inventoryRepository.AddPart(ctx, part)
}
