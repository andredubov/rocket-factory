package inventory

import "context"

func (i *inventoryService) DeletePart(ctx context.Context, uuid string) error {
	return i.inventoryRepository.DeletePart(ctx, uuid)
}
