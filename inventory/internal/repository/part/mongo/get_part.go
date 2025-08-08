package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

func (r *InventoryRepository) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	// Создаем фильтр для поиска по UUID
	filter := bson.M{"_id": uuid}

	repoPart := repoModel.Part{}

	err := r.Collection.FindOne(ctx, filter).Decode(&repoPart)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repository.ErrPartWithUUIDNotFound(uuid)
		}
		return nil, err
	}

	// Конвертируем репозиторную модель в доменную
	part := converter.PartToModel(repoPart)

	return &part, nil
}
