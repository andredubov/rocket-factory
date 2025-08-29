package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// TestOrderUpdateInfoToRepoModel_PartUUIDs verifies correct handling of PartUUIDs in update info
func TestOrderUpdateInfoToRepoModel_PartUUIDs(t *testing.T) {
	// Arrange
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	updateInfo := model.OrderUpdateInfo{
		OrderUUID: uuid.New(),
		PartUUIDs: partUUIDs,
	}

	// Act
	result := converter.OrderUpdateInfoToRepoModel(updateInfo)

	// Assert
	require.Equal(t, updateInfo.OrderUUID, result.OrderUUID)
	require.Len(t, result.PartUUIDs, len(partUUIDs))
	require.Equal(t, partUUIDs, result.PartUUIDs)
}

// TestOrderUpdateInfoToRepoModel_NilPartUUIDs verifies that nil PartUUIDs are handled correctly
func TestOrderUpdateInfoToRepoModel_NilPartUUIDs(t *testing.T) {
	// Arrange
	updateInfo := model.OrderUpdateInfo{
		OrderUUID: uuid.New(),
		PartUUIDs: nil, // Explicitly nil
	}

	// Act
	result := converter.OrderUpdateInfoToRepoModel(updateInfo)

	// Assert
	require.Equal(t, updateInfo.OrderUUID, result.OrderUUID)
	require.Nil(t, result.PartUUIDs)
}

// TestOrderUpdateInfoToRepoModel_EmptyPartUUIDs verifies that empty PartUUIDs slice is handled correctly
func TestOrderUpdateInfoToRepoModel_EmptyPartUUIDs(t *testing.T) {
	// Arrange
	updateInfo := model.OrderUpdateInfo{
		OrderUUID: uuid.New(),
		PartUUIDs: []uuid.UUID{}, // Empty slice
	}

	// Act
	result := converter.OrderUpdateInfoToRepoModel(updateInfo)

	// Assert
	require.Equal(t, updateInfo.OrderUUID, result.OrderUUID)
	require.Empty(t, result.PartUUIDs)
}

// TestOrderUpdateInfoToRepoModel_PartUUIDsWithOtherFields verifies PartUUIDs handling when other fields are present
func TestOrderUpdateInfoToRepoModel_PartUUIDsWithOtherFields(t *testing.T) {
	// Arrange
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New()}
	userUUID := uuid.New()
	totalPrice := 150.0
	status := model.OrderStatusPaid

	updateInfo := model.OrderUpdateInfo{
		OrderUUID:  uuid.New(),
		UserUUID:   &userUUID,
		TotalPrice: &totalPrice,
		Status:     &status,
		PartUUIDs:  partUUIDs,
	}

	// Act
	result := converter.OrderUpdateInfoToRepoModel(updateInfo)

	// Assert
	require.Equal(t, updateInfo.OrderUUID, result.OrderUUID)
	require.Equal(t, userUUID, result.UserUUID)
	require.Equal(t, totalPrice, result.TotalPrice)
	require.Equal(t, repoModel.OrderStatus(status), result.Status)
	require.Len(t, result.PartUUIDs, len(partUUIDs))
	require.Equal(t, partUUIDs, result.PartUUIDs)
}
