package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
	mongoRepository "github.com/andredubov/rocket-factory/inventory/internal/repository/part/mongo"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/mongo/mocks"
)

func TestUpdatePart_Success(t *testing.T) {
	inputPart := model.Part{
		Uuid: "123",
		Name: "Valid Part",
	}

	// Setup mock
	mockCollection := new(mocks.MongoCollection)

	// Создаем SingleResult, который не вернет ошибку
	// Для этого можно использовать реальный FindOne с пустым результатом
	// или создать SingleResult с помощью bson.Raw
	rawBytes, _ := bson.Marshal(bson.M{"_id": "123", "name": "Valid Part"})
	raw := bson.Raw(rawBytes)
	singleResult := mongo.NewSingleResultFromDocument(raw, nil, nil)

	mockCollection.On("FindOneAndUpdate",
		mock.Anything,
		bson.M{"_id": "123"},
		mock.Anything,
		mock.Anything,
	).Return(singleResult)

	// Create repository with mock
	repo := &mongoRepository.InventoryRepository{
		Collection: mockCollection,
	}

	// Execute function
	err := repo.UpdatePart(context.Background(), inputPart)

	// Assertions
	assert.NoError(t, err)
	mockCollection.AssertExpectations(t)
}

func TestUpdatePart_PartNotFount(t *testing.T) {
	var (
		// Helper function to create a SingleResult with error
		newSingleResultWithErr = func(_ error) *mongo.SingleResult {
			return &mongo.SingleResult{
				// SingleResult internals aren't exported, but we can use Decode to simulate behavior
			}
		}
		inputPart = model.Part{
			Uuid: "456",
			Name: "Non-existent Part",
		}
		expectedError = repository.ErrPartWithUUIDNotFound("456")
	)

	// Setup mock
	mockCollection := new(mocks.MongoCollection)
	mockCollection.On("FindOneAndUpdate",
		mock.Anything,
		bson.M{"_id": "456"},
		mock.Anything,
		mock.Anything,
	).Return(newSingleResultWithErr(mongo.ErrNoDocuments))

	// Create repository with mock
	repo := &mongoRepository.InventoryRepository{
		Collection: mockCollection,
	}

	// Execute function
	err := repo.UpdatePart(context.Background(), inputPart)

	// Assertions
	if expectedError != nil {
		assert.Error(t, err)
		if errors.Is(expectedError, repository.ErrPartWithUUIDNotFound("456")) {
			assert.Equal(t, expectedError.Error(), err.Error())
		} else {
			assert.Equal(t, expectedError, err)
		}
	} else {
		assert.NoError(t, err)
	}

	// Verify mock expectations
	mockCollection.AssertExpectations(t)
}

func TestUpdatePart_DatabaseError(t *testing.T) {
	var (
		inputPart = model.Part{
			Uuid: "789",
			Name: "Part with Error",
		}
		dbError = errors.New("database connection failed")
	)

	// Setup mock
	mockCollection := new(mocks.MongoCollection)

	doc := converter.PartToRepoModel(inputPart)
	raw, err := bson.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем SingleResult с ошибкой
	sr := mongo.NewSingleResultFromDocument(
		bson.Raw(raw),
		dbError, // Здесь устанавливаем нашу ошибку
		bson.NewRegistry(),
	)

	mockCollection.On("FindOneAndUpdate",
		mock.Anything,
		bson.M{"_id": "789"},
		mock.Anything,
		mock.Anything,
	).Return(sr)

	// Create repository with mock
	repo := &mongoRepository.InventoryRepository{
		Collection: mockCollection,
	}

	// Execute function
	err = repo.UpdatePart(context.Background(), inputPart)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, dbError, err) // Проверяем, что возвращается именно та ошибка, которую вернула база

	// Verify mock expectations
	mockCollection.AssertExpectations(t)
}
