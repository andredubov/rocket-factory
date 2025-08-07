package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository"
	"github.com/andredubov/rocket-factory/order/internal/repository/order/postgres"
)

// Helper function to capture log output
func captureLogOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func TestGetOrderByOrderUUID_Success(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)

	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()
	paymentMethod := string(model.PaymentMethodCard)
	orderStatus := string(model.OrderStatusPaid)
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)

	columns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	// Expectations
	mock.ExpectBegin()

	// Expect order details query
	mock.ExpectQuery(`SELECT ` + strings.Join(columns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(columns).
				AddRow(
					orderUUID,
					userUUID,
					100.50,
					&transactionUUID, // as *uuid.UUID (pointer)
					&paymentMethod,   // as string
					orderStatus,      // as string
					createdAt,        // as time.Time
					&updatedAt,       // as *time.Time
				),
		)

	// Expect parts query
	mock.ExpectQuery(`SELECT ` + postgres.PartUUIDTableColumn + ` FROM ` + postgres.OrderPartsTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(
			pgxmock.NewRows([]string{postgres.PartUUIDTableColumn}).AddRow(partUUID1).AddRow(partUUID2),
		)

	mock.ExpectCommit()

	// Test
	order, err := repo.GetOrder(context.Background(), orderUUID)

	// Verify
	require.NoError(t, err, "GetOrder should not return an error")
	require.NotNil(t, order, "Order should not be nil")
	require.Equal(t, orderUUID, order.OrderUUID)
	require.Equal(t, userUUID, order.UserUUID)
	require.Equal(t, 100.50, order.TotalPrice)
	require.Equal(t, model.OrderStatusPaid, order.Status)
	require.Len(t, order.PartUUIDs, 2)
	require.Contains(t, order.PartUUIDs, partUUID1)
	require.Contains(t, order.PartUUIDs, partUUID2)
	require.NotNil(t, order.PaymentInfo)
	require.Equal(t, transactionUUID, order.PaymentInfo.TransactionUUID)
	require.Equal(t, model.PaymentMethodCard, order.PaymentInfo.PaymentMethod)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderByOrderUUID_WithoutPaymentInfo(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()
	userUUID := uuid.New()
	createdAt := time.Now()
	orderStatus := string(model.OrderStatusPending)

	columns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	// Expectations
	mock.ExpectBegin()

	// Expect order details query with NULL payment info
	mock.ExpectQuery(`SELECT ` + strings.Join(columns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(pgxmock.NewRows(columns).AddRow(
			orderUUID, userUUID, 50.25, nil, nil, orderStatus, createdAt, nil,
		))

	// Expect empty parts
	mock.ExpectQuery(`SELECT ` + postgres.PartUUIDTableColumn + ` FROM ` + postgres.OrderPartsTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(pgxmock.NewRows([]string{postgres.PartUUIDTableColumn}))

	mock.ExpectCommit()

	// Test
	order, err := repo.GetOrder(context.Background(), orderUUID)

	// Verify
	require.NoError(t, err, "GetOrder should not return an error")
	require.NotNil(t, order, "Order should not be nil")
	require.Equal(t, orderUUID, order.OrderUUID)
	require.Equal(t, userUUID, order.UserUUID)
	require.Equal(t, 50.25, order.TotalPrice)
	require.Equal(t, model.OrderStatusPending, order.Status)
	require.Empty(t, order.PartUUIDs)
	require.Nil(t, order.PaymentInfo)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrder_NotFound(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()

	columns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	// Expectations
	mock.ExpectBegin()

	// Expect no rows for order details
	mock.ExpectQuery(`SELECT ` + strings.Join(columns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnError(pgx.ErrNoRows)

	mock.ExpectRollback()

	// Test
	order, err := repo.GetOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "GetOrder should return an error")
	require.Nil(t, order, "Order should be nil")
	require.EqualError(t, err, repository.ErrOrderNotFoundWith(orderUUID).Error())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderByOrderUUID_BeginTxError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()

	// Expectations
	expectedErr := errors.New("begin error")
	mock.ExpectBegin().WillReturnError(expectedErr)

	// Test
	order, err := repo.GetOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "GetOrder should return an error")
	require.Nil(t, order, "Order should be nil")
	require.Contains(t, err.Error(), "failed to begin transaction")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderByOrderUUID_DetailsQueryError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()

	columns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	// Expectations
	mock.ExpectBegin()

	// Expect query error
	expectedErr := errors.New("query error")
	mock.ExpectQuery(`SELECT ` + strings.Join(columns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnError(expectedErr)

	mock.ExpectRollback()

	// Test
	order, err := repo.GetOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "GetOrder should return an error")
	require.Nil(t, order, "Order should be nil")
	require.Contains(t, err.Error(), "failed to execute query")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderByOrderUUID_PartsQueryError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()
	userUUID := uuid.New()
	createdAt := time.Now()
	orderStatus := string(model.OrderStatusPending)

	columns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	// Expectations
	mock.ExpectBegin()

	// Expect order details
	mock.ExpectQuery(`SELECT ` + strings.Join(columns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(columns).
				AddRow(orderUUID, userUUID, 75.0, nil, nil, orderStatus, createdAt, nil),
		)

	// Expect parts query error
	expectedErr := errors.New("parts query error")
	mock.ExpectQuery(`SELECT ` + postgres.PartUUIDTableColumn + ` FROM ` + postgres.OrderPartsTable).
		WithArgs(orderUUID.String()).
		WillReturnError(expectedErr)

	mock.ExpectRollback()

	// Test
	order, err := repo.GetOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "GetOrder should return an error")
	require.Nil(t, order, "Order should be nil")
	require.Contains(t, err.Error(), "failed to execute parts query")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderByOrderUUID_InvalidStatus(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()
	userUUID := uuid.New()
	createdAt := time.Now()
	orderStatus := "INVALID_STATUS"

	columns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	// Expectations
	mock.ExpectBegin()

	// Expect order with invalid status
	mock.ExpectQuery(`SELECT ` + strings.Join(columns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(columns).AddRow(orderUUID, userUUID, 75.0, nil, nil, orderStatus, createdAt, nil),
		)

	mock.ExpectRollback()

	// Test
	order, err := repo.GetOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "GetOrder should return an error")
	require.Nil(t, order, "Order should be nil")
	require.Contains(t, err.Error(), "invalid order status")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrder_InvalidPaymentMethod(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()
	createdAt := time.Now()
	orderStatus := string(model.OrderStatusPaid)
	paymentMethod := "INVALID_METHOD"

	columns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	// Expectations
	mock.ExpectBegin()

	// Expect order with invalid payment method
	mock.ExpectQuery(`SELECT ` + strings.Join(columns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(columns).AddRow(orderUUID, userUUID, 100.0, &transactionUUID, &paymentMethod, orderStatus, createdAt, nil),
		)

	mock.ExpectRollback()

	// Test
	order, err := repo.GetOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "GetOrder should return an error")
	require.Nil(t, order, "Order should be nil")
	require.Contains(t, err.Error(), "invalid payment method")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrder_CommitError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()
	createdAt := time.Now()
	orderStatus := string(model.OrderStatusPending)
	partUUID := uuid.New()

	columns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	// Expectations
	mock.ExpectBegin()

	// Expect order details
	mock.ExpectQuery(`SELECT ` + strings.Join(columns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(columns).AddRow(orderUUID, userUUID, 100.0, &transactionUUID, nil, orderStatus, createdAt, nil),
		)

	// Expect parts
	mock.ExpectQuery(`SELECT ` + postgres.PartUUIDTableColumn + ` FROM ` + postgres.OrderPartsTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(
			pgxmock.NewRows([]string{postgres.PartUUIDTableColumn}).AddRow(partUUID),
		)

	// Expect commit error
	expectedErr := errors.New("commit error")
	mock.ExpectCommit().WillReturnError(expectedErr)

	// Test
	order, err := repo.GetOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "GetOrder should return an error")
	require.Nil(t, order, "Order should be nil")
	require.Contains(t, err.Error(), "failed to commit transaction")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderByOrderUUID_RollbackOnError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)

	orderUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()

	// Simulate error in order details query
	mock.ExpectQuery(`SELECT .* FROM orders`).
		WithArgs(orderUUID.String()).
		WillReturnError(fmt.Errorf("database error"))

	// We expect rollback to be called
	mock.ExpectRollback()

	// Test
	order, err := repo.GetOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "Should return error")
	require.Nil(t, order, "Should not return order")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetOrderByOrderUUID_RollbackFailure(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)

	orderUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()

	// Simulate error in order details query
	mock.ExpectQuery(`SELECT .* FROM orders`).
		WithArgs(orderUUID.String()).
		WillReturnError(fmt.Errorf("database error"))

	// Simulate rollback failure
	mock.ExpectRollback().
		WillReturnError(fmt.Errorf("rollback failed"))

	// Capture log output
	logOutput := captureLogOutput(func() {
		// Test
		order, err := repo.GetOrder(context.Background(), orderUUID)

		// Verify
		require.Error(t, err, "Should return error")
		require.Nil(t, order, "Should not return order")
	})

	// Verify rollback error was logged
	require.Contains(t, logOutput, "failed to rollback transaction", "Should log rollback error")
	require.NoError(t, mock.ExpectationsWereMet(), "Should meet all expectations")
}

func TestGetOrderByOrderUUID_PartScanError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()
	userUUID := uuid.New()
	createdAt := time.Now()
	orderStatus := string(model.OrderStatusPending)

	columns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	// Expectations
	mock.ExpectBegin()

	// Expect order details
	mock.ExpectQuery(`SELECT ` + strings.Join(columns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(columns).
				AddRow(orderUUID, userUUID, 75.0, nil, nil, orderStatus, createdAt, nil),
		)

	// Expect parts query to return a row with invalid data that will cause scan to fail
	// We'll use a string that can't be parsed as UUID
	invalidUUID := "not-a-valid-uuid"
	mock.ExpectQuery(`SELECT ` + postgres.PartUUIDTableColumn + ` FROM ` + postgres.OrderPartsTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(
			pgxmock.NewRows([]string{postgres.PartUUIDTableColumn}).AddRow(invalidUUID),
		)

	mock.ExpectRollback()

	// Test
	order, err := repo.GetOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "GetOrder should return an error")
	require.Nil(t, order, "Order should be nil")
	require.Contains(t, err.Error(), "failed to scan part UUID", "Error should mention scan failure")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrder_PartsRowsError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()
	userUUID := uuid.New()
	createdAt := time.Now()
	orderStatus := string(model.OrderStatusPending)

	columns := []string{
		postgres.UUIDTableColumn,
		postgres.UserUUIDTableColumn,
		postgres.TotalPriceTableColumn,
		postgres.TransactionUUIDTableColumn,
		postgres.PaymentMethodTableColumn,
		postgres.StatusTableColumn,
		postgres.CreatedAtTableColumn,
		postgres.UpdatedAtTableColumn,
	}

	// Expectations
	mock.ExpectBegin()

	// Expect order details
	mock.ExpectQuery(`SELECT ` + strings.Join(columns, ", ") + ` FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(
			pgxmock.NewRows(columns).
				AddRow(orderUUID, userUUID, 75.0, nil, nil, orderStatus, createdAt, nil),
		)

	// Simulate error during rows iteration
	expectedErr := errors.New("rows iteration error")
	mock.ExpectQuery(`SELECT ` + postgres.PartUUIDTableColumn + ` FROM ` + postgres.OrderPartsTable).
		WithArgs(orderUUID.String()).
		WillReturnRows(
			pgxmock.NewRows([]string{postgres.PartUUIDTableColumn}).
				AddRow(uuid.New()).       // First row is fine
				RowError(1, expectedErr), // Error on second row
		)

	mock.ExpectRollback()

	// Test
	order, err := repo.GetOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "GetOrder should return an error")
	require.Nil(t, order, "Order should be nil")
	require.Contains(t, err.Error(), "error processing results", "Error should mention results processing")
	require.ErrorIs(t, err, expectedErr, "Error should wrap the original rows error")
	require.NoError(t, mock.ExpectationsWereMet())
}
