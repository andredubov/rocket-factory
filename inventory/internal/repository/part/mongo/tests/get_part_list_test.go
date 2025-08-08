package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/dvln/testify/require"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	mongodb "github.com/andredubov/rocket-factory/inventory/internal/repository/part/mongo"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/mongo/mocks"
)

func TestGetPartList_DatabaseError(t *testing.T) {
	var (
		filter      = model.PartFilter{}
		expectedErr = errors.New("some mongo database error")
	)

	// Setup mocks
	mockCollection := mocks.NewMongoCollection(t)
	mockCursor := mocks.NewMongoCursor(t)

	mockCollection.On("Find", mock.Anything, bson.M{}, mock.Anything).Return(nil, expectedErr)

	// Create repository with mock collection
	repo := &mongodb.InventoryRepository{
		Collection: mockCollection,
	}

	// Execute
	parts, err := repo.GetPartList(context.Background(), filter)

	// Assertions
	// assert.Equal(t, expectedErr, err)
	require.Error(t, err)
	assert.Equal(t, err, expectedErr)
	assert.Nil(t, parts)

	// Verify mock expectations
	mockCollection.AssertExpectations(t)
	mockCursor.AssertExpectations(t)
}

func TestGetPartList_NoDocumentsFound(t *testing.T) {
	var (
		filter      = model.PartFilter{}
		expectedErr = mongo.ErrNoDocuments
	)

	// Setup mocks
	mockCollection := mocks.NewMongoCollection(t)
	mockCursor := mocks.NewMongoCursor(t)

	mockCollection.On("Find", mock.Anything, bson.M{}, mock.Anything).Return(nil, expectedErr)

	// Create repository with mock collection
	repo := &mongodb.InventoryRepository{
		Collection: mockCollection,
	}

	// Execute
	parts, err := repo.GetPartList(context.Background(), filter)

	// Assertions
	// assert.Equal(t, expectedErr, err)
	require.Error(t, err)
	assert.True(t, errors.Is(err, model.ErrPartNotFound))
	assert.Nil(t, parts)

	// Verify mock expectations
	mockCollection.AssertExpectations(t)
	mockCursor.AssertExpectations(t)
}
