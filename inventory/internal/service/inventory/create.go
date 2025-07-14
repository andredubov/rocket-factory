package inventory

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
)

func (i *inventoryService) AddPart(ctx context.Context, part model.Part) error {
	repoPart := converter.PartToRepoModel(part)
	return i.inventoryRepository.AddPart(ctx, repoPart)
}
