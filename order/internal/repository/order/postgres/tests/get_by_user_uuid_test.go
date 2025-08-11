package tests

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/order/postgres"
)

func TestGetUserOrders_Success(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)

	userUUID := uuid.New()
	orderUUID1 := uuid.New()
	orderUUID2 := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()
	transactionUUID := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)
	paymentMethod := string(model.PaymentMethodCard)
	orderStatus := string(model.OrderStatusPaid)

	// Expectations
	mock.ExpectBegin()

	// Expect orders query
	orderColumns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	mock.ExpectQuery(`SELECT ` + strings.Join(orderColumns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(userUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(orderColumns).
				AddRow(
					orderUUID1,
					userUUID,
					100.50,
					&transactionUUID,
					&paymentMethod,
					orderStatus,
					createdAt,
					&updatedAt,
				).
				AddRow(
					orderUUID2,
					userUUID,
					200.75,
					nil,
					nil,
					orderStatus,
					createdAt.Add(-time.Hour),
					nil,
				),
		)

	// For parts query
	mock.ExpectQuery(`SELECT `+postgres.OrderUUIDTableColumn+`, `+postgres.PartUUIDTableColumn+` FROM `+
		postgres.OrderPartsTable+` WHERE `+postgres.OrderUUIDTableColumn+` IN \(\$1,\$2\)`).
		WithArgs(orderUUID1, orderUUID2).
		WillReturnRows(
			pgxmock.NewRows([]string{postgres.OrderUUIDTableColumn, postgres.PartUUIDTableColumn}).
				AddRow(orderUUID1, partUUID1).
				AddRow(orderUUID1, partUUID2),
		)

	mock.ExpectCommit()

	// Test
	orders, err := repo.GetUserOrders(context.Background(), userUUID)

	// Verify
	require.NoError(t, err, "Should not return error")
	require.Len(t, orders, 2, "Should return 2 orders")

	// Verify first order
	require.Equal(t, orderUUID1, orders[0].OrderUUID)
	require.Equal(t, userUUID, orders[0].UserUUID)
	require.Equal(t, 100.50, orders[0].TotalPrice)
	require.Equal(t, model.OrderStatusPaid, orders[0].Status)
	require.Len(t, orders[0].PartUUIDs, 2)
	require.Contains(t, orders[0].PartUUIDs, partUUID1)
	require.Contains(t, orders[0].PartUUIDs, partUUID2)
	require.NotNil(t, orders[0].PaymentInfo)
	require.Equal(t, transactionUUID, orders[0].PaymentInfo.TransactionUUID)
	require.Equal(t, model.PaymentMethodCard, orders[0].PaymentInfo.PaymentMethod)

	// Verify second order
	require.Equal(t, orderUUID2, orders[1].OrderUUID)
	require.Equal(t, userUUID, orders[1].UserUUID)
	require.Equal(t, 200.75, orders[1].TotalPrice)
	require.Equal(t, model.OrderStatusPaid, orders[1].Status)
	require.Empty(t, orders[1].PartUUIDs)
	require.Nil(t, orders[1].PaymentInfo)

	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetUserOrders_NoOrders(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	userUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()

	// Expect orders query to return no rows
	orderColumns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	mock.ExpectQuery(`SELECT ` + strings.Join(orderColumns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(userUUID.String()).
		WillReturnRows(pgxmock.NewRows(orderColumns))

	// No parts query expected since no orders
	mock.ExpectCommit()

	// Test
	orders, err := repo.GetUserOrders(context.Background(), userUUID)

	// Verify
	require.NoError(t, err, "Should not return error")
	require.Empty(t, orders, "Should return empty slice")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetUserOrders_OrderQueryError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	userUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()

	// Simulate error in orders query
	mock.ExpectQuery(`SELECT .* FROM ` + postgres.OrdersTable).
		WithArgs(userUUID.String()).
		WillReturnError(fmt.Errorf("database error"))

	// Expect rollback
	mock.ExpectRollback()

	// Test
	orders, err := repo.GetUserOrders(context.Background(), userUUID)

	// Verify
	require.Error(t, err, "Should return error")
	require.Nil(t, orders, "Should not return orders")
	require.Contains(t, err.Error(), "failed to query user orders", "Error should indicate query failure")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetUserOrders_PartsQueryError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	userUUID := uuid.New()
	orderUUID := uuid.New()
	transactionUUID := uuid.New()
	createdAt := time.Now()
	paymentMethod := string(model.PaymentMethodCard)
	orderStatus := string(model.OrderStatusPaid)

	// Expectations
	mock.ExpectBegin()

	// Expect orders query
	orderColumns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	mock.ExpectQuery(`SELECT ` + strings.Join(orderColumns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(userUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(orderColumns).
				AddRow(
					orderUUID,
					userUUID,
					100.50,
					&transactionUUID,
					&paymentMethod,
					orderStatus,
					createdAt,
					nil,
				),
		)

	// Simulate error in parts query
	mock.ExpectQuery(`SELECT .* FROM ` + postgres.OrderPartsTable).
		WithArgs(orderUUID).
		WillReturnError(fmt.Errorf("parts database error"))

	// Expect rollback
	mock.ExpectRollback()

	// Test
	orders, err := repo.GetUserOrders(context.Background(), userUUID)

	// Verify
	require.Error(t, err, "Should return error")
	require.Nil(t, orders, "Should not return orders")
	require.Contains(t, err.Error(), "failed to query order parts", "Error should indicate parts query failure")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetUserOrders_CommitError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	userUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()

	// Expect empty orders result
	orderColumns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	mock.ExpectQuery(`SELECT ` + strings.Join(orderColumns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(userUUID.String()).
		WillReturnRows(pgxmock.NewRows(orderColumns))

	// Simulate commit error
	mock.ExpectCommit().
		WillReturnError(fmt.Errorf("commit failed"))

	// Capture log output
	var logOutput string
	log.SetOutput(&testLogWriter{&logOutput})
	defer log.SetOutput(nil)

	// Test
	orders, err := repo.GetUserOrders(context.Background(), userUUID)

	// Verify
	require.Error(t, err, "Should return error")
	require.Nil(t, orders, "Should not return orders")
	require.Contains(t, err.Error(), "failed to commit transaction", "Error should indicate commit failure")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetUserOrders_WithInvalidStatus(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	userUUID := uuid.New()
	orderUUID := uuid.New()
	createdAt := time.Now()

	// Expectations
	mock.ExpectBegin()

	// Expect orders query with invalid status
	orderColumns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	mock.ExpectQuery(`SELECT ` + strings.Join(orderColumns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(userUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(orderColumns).
				AddRow(
					orderUUID,
					userUUID,
					100.50,
					nil,
					nil,
					"INVALID_STATUS", // Invalid status
					createdAt,
					nil,
				),
		)

	// Expect rollback due to invalid status
	mock.ExpectRollback()

	// Test
	orders, err := repo.GetUserOrders(context.Background(), userUUID)

	// Verify
	require.Error(t, err, "Should return error")
	require.Nil(t, orders, "Should not return orders")
	require.Contains(t, err.Error(), "invalid order status", "Error should indicate status validation failure")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetUserOrders_BeginTransactionError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	userUUID := uuid.New()

	// Simulate error when beginning transaction
	mock.ExpectBegin().
		WillReturnError(fmt.Errorf("cannot begin transaction"))

	// No other expectations needed since the function should fail immediately

	// Test
	orders, err := repo.GetUserOrders(context.Background(), userUUID)

	// Verify
	require.Error(t, err, "Should return error")
	require.Nil(t, orders, "Should not return orders")
	require.Contains(t, err.Error(), "failed to begin transaction", "Error should indicate transaction begin failure")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetUserOrders_RollbackError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	userUUID := uuid.New()

	// Capture log output
	logOutput := captureLogOutput(func() {
		// Expectations
		mock.ExpectBegin()

		// Simulate error in orders query
		mock.ExpectQuery(`SELECT .* FROM ` + postgres.OrdersTable).
			WithArgs(userUUID.String()).
			WillReturnError(fmt.Errorf("database error"))

		// Simulate rollback error
		mock.ExpectRollback().
			WillReturnError(fmt.Errorf("rollback failed"))

		// Test
		orders, err := repo.GetUserOrders(context.Background(), userUUID)

		// Verify
		require.Error(t, err, "Should return error")
		require.Nil(t, orders, "Should not return orders")
		require.Contains(t, err.Error(), "failed to query user orders", "Should return original error")
	})

	// Verify rollback error was logged
	require.Contains(t, logOutput, "failed to rollback transaction", "Should log rollback error")
	require.Contains(t, logOutput, "rollback failed", "Should log rollback error details")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetUserOrders_RowsIterationError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	userUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()

	// Expect orders query
	orderColumns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	// Create rows that will return error during iteration
	rows := pgxmock.NewRows(orderColumns).
		AddRow(
			uuid.New(),
			userUUID,
			100.50,
			nil,
			nil,
			"PAID",
			time.Now(),
			nil,
		).RowError(1, fmt.Errorf("iteration error")) // Add error after first row

	mock.ExpectQuery(`SELECT ` + strings.Join(orderColumns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(userUUID.String()).
		WillReturnRows(rows)

	// Expect rollback due to error
	mock.ExpectRollback()

	// Test
	orders, err := repo.GetUserOrders(context.Background(), userUUID)

	// Verify
	require.Error(t, err, "Should return error")
	require.Nil(t, orders, "Should not return orders")
	require.Contains(t, err.Error(), "error during orders iteration", "Error should indicate iteration failure")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetUserOrders_PartsIterationError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	userUUID := uuid.New()
	orderUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()

	// Expect orders query to return one order
	orderColumns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	mock.ExpectQuery(`SELECT ` + strings.Join(orderColumns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(userUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(orderColumns).
				AddRow(
					orderUUID,
					userUUID,
					100.50,
					nil,
					nil,
					"PAID",
					time.Now(),
					nil,
				),
		)

	// Create parts rows that will return error during iteration
	partsRows := pgxmock.NewRows([]string{postgres.OrderUUIDTableColumn, postgres.PartUUIDTableColumn}).
		AddRow(orderUUID, uuid.New()).
		RowError(1, fmt.Errorf("parts iteration error")) // Add error after first row

	mock.ExpectQuery(`SELECT ` + postgres.OrderUUIDTableColumn + `, ` + postgres.PartUUIDTableColumn + ` FROM ` + postgres.OrderPartsTable).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(partsRows)

	// Expect rollback due to error
	mock.ExpectRollback()

	// Test
	orders, err := repo.GetUserOrders(context.Background(), userUUID)

	// Verify
	require.Error(t, err, "Should return error")
	require.Nil(t, orders, "Should not return orders")
	require.Contains(t, err.Error(), "error during parts iteration", "Error should indicate parts iteration failure")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetUserOrders_PartsScanError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	userUUID := uuid.New()
	orderUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()

	// Expect orders query to return one order
	orderColumns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	mock.ExpectQuery(`SELECT ` + strings.Join(orderColumns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(userUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(orderColumns).
				AddRow(
					orderUUID,
					userUUID,
					100.50,
					nil,
					nil,
					"PAID",
					time.Now(),
					nil,
				),
		)

	// Create parts rows with invalid data type to force scan error
	partsRows := pgxmock.NewRows([]string{postgres.OrderUUIDTableColumn, postgres.PartUUIDTableColumn}).
		AddRow("invalid_uuid_format", 12345) // Wrong types to force scan error

	mock.ExpectQuery(`SELECT ` + postgres.OrderUUIDTableColumn + `, ` + postgres.PartUUIDTableColumn + ` FROM ` + postgres.OrderPartsTable).
		WithArgs(pgxmock.AnyArg()).
		WillReturnRows(partsRows)

	// Expect rollback due to error
	mock.ExpectRollback()

	// Test
	orders, err := repo.GetUserOrders(context.Background(), userUUID)

	// Verify
	require.Error(t, err, "Should return error")
	require.Nil(t, orders, "Should not return orders")
	require.Contains(t, err.Error(), "failed to scan part", "Error should indicate scan failure")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetUserOrders_ScanRowError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	userUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()

	// Mock orders query that will return a row with invalid data types
	orderColumns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	// Create a row with invalid data types that will cause scan to fail
	// Using string for UUID field which should be []byte
	mock.ExpectQuery(`SELECT ` + strings.Join(orderColumns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(userUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(orderColumns).
				AddRow(
					"invalid-uuid-format", // Invalid UUID format (should be []byte or uuid.UUID)
					userUUID,
					100.50,
					nil,
					nil,
					"PAID",
					time.Now(),
					nil,
				),
		)

	mock.ExpectRollback()

	// Test
	orders, err := repo.GetUserOrders(context.Background(), userUUID)

	// Verify
	require.Error(t, err, "Should return error")
	require.Nil(t, orders, "Should not return orders")
	require.Contains(t, err.Error(), "failed to scan order", "Error should indicate scan failure")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetUserOrders_InvalidPaymentMethod(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	userUUID := uuid.New()
	orderUUID := uuid.New()
	transactionUUID := uuid.New()
	invalidPaymentMethod := "INVALID_PAYMENT_METHOD"
	createdAt := time.Now()
	orderStatus := string(model.OrderStatusPaid)

	// Expectations
	mock.ExpectBegin()

	// Mock order with invalid payment method
	orderColumns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	mock.ExpectQuery(`SELECT ` + strings.Join(orderColumns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(userUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(orderColumns).
				AddRow(
					orderUUID,
					userUUID,
					100.50,
					&transactionUUID,
					&invalidPaymentMethod,
					orderStatus,
					createdAt,
					nil,
				),
		)

	// We don't expect the parts query to be executed because scanOrderRow will fail first
	// So we expect rollback immediately after the order query

	mock.ExpectRollback()

	// Test
	orders, err := repo.GetUserOrders(context.Background(), userUUID)

	// Verify
	require.Error(t, err, "Should return error for invalid payment method")
	require.Nil(t, orders, "Should not return orders")
	require.Contains(t, err.Error(), "invalid payment method", "Error should indicate payment method validation failure")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}
