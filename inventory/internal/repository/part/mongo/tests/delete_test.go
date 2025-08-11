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
	mongoRepository "github.com/andredubov/rocket-factory/inventory/internal/repository/part/mongo"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/mongo/mocks"
)

func TestInventoryRepository_DeletePart(t *testing.T) {
	tests := []struct {
		name          string
		uuid          string
		mockSetup     func(*mocks.MongoCollection)
		expectedError error
	}{
		{
			name: "successful deletion",
			uuid: "123",
			mockSetup: func(mc *mocks.MongoCollection) {
				mc.On("DeleteOne", mock.Anything, bson.M{"_id": "123"}, mock.Anything).
					Return(&mongo.DeleteResult{DeletedCount: 1}, nil)
			},
			expectedError: nil,
		},
		{
			name: "part not found",
			uuid: "123",
			mockSetup: func(mc *mocks.MongoCollection) {
				mc.On("DeleteOne", mock.Anything, bson.M{"_id": "123"}, mock.Anything).
					Return(&mongo.DeleteResult{DeletedCount: 0}, nil)
			},
			expectedError: model.ErrPartWithUUIDNotFound("123"),
		},
		{
			name: "database error",
			uuid: "123",
			mockSetup: func(mc *mocks.MongoCollection) {
				mc.On("DeleteOne", mock.Anything, bson.M{"_id": "123"}, mock.Anything).
					Return(nil, errors.New("connection error"))
			},
			expectedError: errors.New("connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок коллекции
			mockCollection := new(mocks.MongoCollection)

			// Настраиваем ожидания мока
			tt.mockSetup(mockCollection)

			// Создаем репозиторий с моком
			repo := &mongoRepository.InventoryRepository{
				Collection: mockCollection,
			}

			// Вызываем тестируемый метод
			err := repo.DeletePart(context.Background(), tt.uuid)

			// Проверяем результаты
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверяем, что все ожидания мока были выполнены
			mockCollection.AssertExpectations(t)
		})
	}
}
