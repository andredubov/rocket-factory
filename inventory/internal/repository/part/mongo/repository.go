package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/andredubov/rocket-factory/inventory/internal/service"
)

const (
	PartsCollection = "parts"
)

type inventoryRepository struct {
	collection *mongo.Collection
}

// NewInventoryRepository creates a mongodb inventory repository instance
func NewInventoryRepository(db *mongo.Database) service.InventoryRepository {
	collection := db.Collection(PartsCollection)

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	return &inventoryRepository{
		collection: collection,
	}
}
