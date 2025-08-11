package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
)

func (r *InventoryRepository) AddPart(ctx context.Context, part model.Part) error {
	repoPart := converter.PartToRepoModel(part)

	// Вставка документа. MongoDB автоматически проверяет уникальность _id (Uuid).
	_, err := r.Collection.InsertOne(ctx, repoPart)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return model.ErrPartWithUUIDExists(part.Uuid)
		}
		return err
	}

	return nil
}
