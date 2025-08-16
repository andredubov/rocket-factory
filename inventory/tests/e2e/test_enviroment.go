//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	gofakeit.Seed(time.Now().UnixNano()) //nolint:errcheck,gosec

	parsed := uuid.New()
	uuid := parsed.String()
	now := time.Now()

	document := bson.M{
		"_id":            primitive.Binary{Subtype: 0x04, Data: parsed[:]},
		"name":           gofakeit.RandomString([]string{"Main Engine", "Thruster", "Fuel Tank", "Left Wing", "Right Wing"}),
		"description":    gofakeit.Sentence(8),
		"price":          int64(gofakeit.Price(100, 300_000)),
		"stock_quantity": int64(gofakeit.Number(1, 25)),
		"category":       int32(gofakeit.Number(0, 4)), // nolint:gosec
		"dimensions": bson.M{
			"length": gofakeit.Float64Range(1, 1000),
			"width":  gofakeit.Float64Range(1, 1000),
			"height": gofakeit.Float64Range(1, 1000),
			"weight": gofakeit.Float64Range(1, 1000),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": gofakeit.URL(),
		},
		"tags":       []string{gofakeit.Word(), gofakeit.Word()},
		"metadata":   bson.M{"key": gofakeit.Word()},
		"created_at": primitive.NewDateTimeFromTime(now),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	databaseName := mongoDBName()
	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, document)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
