package mongodb

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

func (r *InventoryRepository) GetPartList(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	repoFilter := converter.PartFilterToRepoModel(filter)

	// Строим MongoDB-фильтр
	mongoFilter := buildMongoFilter(repoFilter)

	opts := options.Find()

	cursor, err := r.Collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}
		return nil, err
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Println("failed to close cursor: ", err)
		}
	}()

	var parts []repoModel.Part

	if err = cursor.All(ctx, &parts); err != nil {
		return nil, err
	}

	if len(parts) == 0 {
		return nil, model.ErrPartNotFound
	}

	return converter.PartsToModel(parts), nil
}

// buildMongoFilter конструирует фильтр для MongoDB на основе критериев
func buildMongoFilter(filter repoModel.PartFilter) bson.M {
	mongoFilter := bson.M{}

	// OR-условия для каждого поля
	if len(filter.UUIDs) > 0 {
		mongoFilter["_id"] = bson.M{"$in": filter.UUIDs}
	}

	if len(filter.Names) > 0 {
		mongoFilter["name"] = bson.M{"$in": filter.Names}
	}

	if len(filter.Categories) > 0 {
		mongoFilter["category"] = bson.M{"$in": filter.Categories}
	}

	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	return mongoFilter
}
