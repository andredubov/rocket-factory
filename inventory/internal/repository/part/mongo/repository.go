package mongodb

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

type MongoCursor interface {
	Next(ctx context.Context) bool
	TryNext(ctx context.Context) bool
	Decode(val interface{}) error
	Close(ctx context.Context) error
	Err() error
	All(ctx context.Context, results interface{}) error
	ID() int64
	RemainingBatchLength() int
	Current() bson.Raw
	SetBatchSize(int32)
	GetBatchSize() int32
}

type MongoCollection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	FindOneAndUpdate(ctx context.Context, filter, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
}

type InventoryRepository struct {
	Collection MongoCollection
}

// NewInventoryRepository creates a mongodb inventory repository instance
func NewInventoryRepository(ctx context.Context, db *mongo.Database) service.InventoryRepository {
	collection := db.Collection(PartsCollection)

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	return &InventoryRepository{
		Collection: collection,
	}
}
