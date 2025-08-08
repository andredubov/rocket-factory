package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/inventory/internal/repository"
	mongoRepository "github.com/andredubov/rocket-factory/inventory/internal/repository/part/mongo"
	"github.com/andredubov/rocket-factory/inventory/internal/repository/part/mongo/mocks"
)

func TestInventoryRepository_AddPart(t *testing.T) {
	tests := []struct {
		name          string
		part          model.Part
		mockSetup     func(*mocks.MongoCollection)
		expectedError error
	}{
		{
			name: "successful part addition",
			part: model.Part{Uuid: "123"},
			mockSetup: func(mc *mocks.MongoCollection) {
				mc.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).
					Return(&mongo.InsertOneResult{InsertedID: "123"}, nil)
			},
			expectedError: nil,
		},
		{
			name: "duplicate key error",
			part: model.Part{Uuid: "123"},
			mockSetup: func(mc *mocks.MongoCollection) {
				mc.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, mongo.WriteException{
						WriteErrors: []mongo.WriteError{
							{Code: 11000}, // Duplicate key error code
						},
					})
			},
			expectedError: repository.ErrPartWithUUIDExists("123"),
		},
		{
			name: "other error",
			part: model.Part{Uuid: "123"},
			mockSetup: func(mc *mocks.MongoCollection) {
				mc.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).
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
			err := repo.AddPart(context.Background(), tt.part)

			// Проверяем результаты
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверяем, что все ожидания мока были выполнены
			mockCollection.AssertExpectations(t)
		})
	}
}
