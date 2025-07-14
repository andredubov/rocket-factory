package inventory

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
)

func (i *inventoryService) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	repoPart, err := i.inventoryRepository.GetPart(ctx, uuid)
	if err != nil {
		return nil, err
	}

	part := converter.PartToModel(*repoPart)
	return &part, nil
}

func (i *inventoryService) GetPartList(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	repoFilter := converter.PartFilterToRepoModel(filter)
	repoParts, err := i.inventoryRepository.GetPartList(ctx, repoFilter)
	if err != nil {
		return nil, err
	}

	parts := make([]model.Part, 0)
	for _, repoPart := range repoParts {
		parts = append(parts, converter.PartToModel(repoPart))
	}

	return parts, nil
}
