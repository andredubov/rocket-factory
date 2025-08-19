//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	databaseName := mongoDBName()
	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})

	return err
}

func mongoDBName() string {
	if v := os.Getenv("MONGO_INITDB_DATABASE"); v != "" {
		return v
	}

	return "inventory-service-database"
}

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	var (
		uuid     = uuid.New().String()
		now      = time.Now()
		document = repoModel.Part{
			Uuid:          uuid,
			Name:          gofakeit.RandomString([]string{"Main Engine", "Thruster", "Fuel Tank", "Left Wing", "Right Wing"}),
			Description:   "Description",
			Price:         234.56,
			StockQuantity: 13,
			Category:      repoModel.PartCategoryEngine,
			Dimensions: repoModel.Dimensions{
				Length: 34.43,
				Width:  12.45,
				Height: 67.12,
				Weight: 34.68,
			},
			Manufacturer: repoModel.Manufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags:      []string{gofakeit.Word(), gofakeit.Word()},
			Metadata:  map[string]repoModel.Value{},
			CreatedAt: now.Add(-2 * time.Hour),
			UpdatedAt: now.Add(-1 * time.Hour),
		}
	)

	databaseName := mongoDBName()
	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, document)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
