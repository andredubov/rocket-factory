package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
)

func (r *InventoryRepository) DeletePart(ctx context.Context, uuid string) error {
	filter := bson.M{"_id": uuid}

	// Выполняем операцию удаления
	result, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// Проверяем, был ли удален документ
	if result.DeletedCount == 0 {
		return model.ErrPartWithUUIDNotFound(uuid)
	}

	return nil
}
