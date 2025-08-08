package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
)

func (r *InventoryRepository) UpdatePart(ctx context.Context, part model.Part) error {
	repoPart := converter.PartToRepoModel(part)

	filter := bson.M{"_id": part.Uuid}
	update := bson.M{
		"$set":         repoPart,
		"$currentDate": bson.M{"updatedAt": true},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := r.Collection.FindOneAndUpdate(ctx, filter, update, opts).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return repository.ErrPartWithUUIDNotFound(part.Uuid)
		}
		return err
	}

	return nil
}
