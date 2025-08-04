package inventory

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
)

func (i *inventoryService) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := i.inventoryRepository.GetPart(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return part, nil
}

func (i *inventoryService) GetPartList(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	parts, err := i.inventoryRepository.GetPartList(ctx, filter)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
