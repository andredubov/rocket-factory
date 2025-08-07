package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/andredubov/rocket-factory/inventory/internal/repository"
)

func (r *inventoryRepository) DeletePart(ctx context.Context, uuid string) error {
	filter := bson.M{"_id": uuid}

	// Выполняем операцию удаления
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// Проверяем, был ли удален документ
	if result.DeletedCount == 0 {
		return repository.ErrPartWithUUIDNotFound(uuid)
	}

	return nil
}
