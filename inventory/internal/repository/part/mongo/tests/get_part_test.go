package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
	mongodb "github.com/andredubov/rocket-factory/inventory/internal/repository/part/mongo"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/mongo/mocks"
)

func stringPtr(s string) *string {
	return &s
}

func TestGetPart_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	uuid := "test-uuid"
	now := time.Now().UTC() // Используем UTC и текущее время для теста

	expectedRepoPart := repoModel.Part{
		Uuid:          uuid,
		Name:          "Test Part",
		Description:   "Test Description",
		Price:         100.50,
		StockQuantity: 10,
		Category:      1, // Engine
		Dimensions: repoModel.Dimensions{
			Length: 10.5,
			Width:  5.5,
			Height: 2.5,
			Weight: 1.2,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "Test Manufacturer",
			Country: "US",
			Website: "https://test.com",
		},
		Tags:      []string{"tag1", "tag2"},
		CreatedAt: now,
		Metadata: map[string]repoModel.Value{
			"key": {
				StringValue: stringPtr("value"),
			},
		},
	}

	expectedModelPart := model.Part{
		Uuid:          uuid,
		Name:          "Test Part",
		Description:   "Test Description",
		Price:         100.50,
		StockQuantity: 10,
		Category:      model.PartCategoryEngine,
		Dimensions: model.Dimensions{
			Length: 10.5,
			Width:  5.5,
			Height: 2.5,
			Weight: 1.2,
		},
		Manufacturer: model.Manufacturer{
			Name:    "Test Manufacturer",
			Country: "US",
			Website: "https://test.com",
		},
		Tags:      []string{"tag1", "tag2"},
		Metadata:  map[string]model.Value{"key": {StringValue: stringPtr("value")}},
		CreatedAt: now,
	}

	mockCollection := mocks.NewMongoCollection(t)
	mockCollection.On("FindOne", ctx, bson.M{"_id": uuid}, mock.Anything).
		Return(mongo.NewSingleResultFromDocument(expectedRepoPart, nil, nil))

	repo := &mongodb.InventoryRepository{
		Collection: mockCollection,
	}

	// Act
	part, err := repo.GetPart(ctx, uuid)

	// Assert
	assert.NoError(t, err)

	// Verify
	assert.Equal(t, expectedModelPart.Uuid, part.Uuid)
	assert.Equal(t, expectedModelPart.Name, part.Name)
	assert.Equal(t, expectedModelPart.Description, part.Description)
	assert.Equal(t, expectedModelPart.Price, part.Price)
	assert.Equal(t, expectedModelPart.StockQuantity, part.StockQuantity)
	assert.Equal(t, expectedModelPart.Category, part.Category)
	assert.Equal(t, expectedModelPart.Dimensions, part.Dimensions)
	assert.Equal(t, expectedModelPart.Manufacturer, part.Manufacturer)
	assert.Equal(t, expectedModelPart.Tags, part.Tags)
	assert.Equal(t, *expectedModelPart.Metadata["key"].StringValue, *part.Metadata["key"].StringValue)
	assert.WithinDuration(t, expectedModelPart.CreatedAt, part.CreatedAt, time.Second) // Допустимая погрешность 1 секунда

	mockCollection.AssertExpectations(t)
}

func TestGetPart_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	uuid := "non-existent-uuid"

	mockCollection := mocks.NewMongoCollection(t)
	mockCollection.On("FindOne", ctx, bson.M{"_id": uuid}, mock.Anything).
		Return(mongo.NewSingleResultFromDocument(repoModel.Part{}, mongo.ErrNoDocuments, nil))

	repo := &mongodb.InventoryRepository{
		Collection: mockCollection,
	}

	// Act
	part, err := repo.GetPart(ctx, uuid)

	// Assert
	assert.Nil(t, part)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, model.ErrPartNotFound))
	mockCollection.AssertExpectations(t)
}

func TestGetPart_DatabaseError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	uuid := "non-existent-uuid"
	dbError := errors.New("database connection failed")

	mockCollection := mocks.NewMongoCollection(t)
	mockCollection.On("FindOne", ctx, bson.M{"_id": uuid}, mock.Anything).
		Return(mongo.NewSingleResultFromDocument(repoModel.Part{}, dbError, nil))

	repo := &mongodb.InventoryRepository{
		Collection: mockCollection,
	}

	// Act
	part, err := repo.GetPart(ctx, uuid)

	// Assert
	assert.Nil(t, part)
	assert.Error(t, err)
	assert.Equal(t, dbError, err)
	mockCollection.AssertExpectations(t)
}
