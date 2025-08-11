package tests

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/repository/order/postgres"
)

func TestDeleteOrder_Success(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	// Test
	err = repo.DeleteOrder(context.Background(), orderUUID)

	// Verify
	require.NoError(t, err, "DeleteOrder should not return an error")
	require.NoError(t, mock.ExpectationsWereMet(), "All expected SQL queries should be executed")
}

func TestDeleteOrder_TransactionBeginError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()

	// Expectations
	mock.ExpectBegin().WillReturnError(errors.New("begin error"))

	// Test
	err = repo.DeleteOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "DeleteOrder should return an error when transaction begin fails")
	require.Contains(t, err.Error(), "failed to begin transaction")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteOrder_ExecError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnError(errors.New("exec error"))
	mock.ExpectRollback()

	// Test
	err = repo.DeleteOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "DeleteOrder should return an error when exec fails")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteOrder_CommitError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit().WillReturnError(errors.New("commit error"))

	// Capture log output
	var logOutput string
	log.SetOutput(&testLogWriter{&logOutput})
	defer log.SetOutput(nil)

	// Test
	err = repo.DeleteOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "DeleteOrder should return an error when commit fails")
	require.Contains(t, err.Error(), "failed to commit transaction")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteOrder_RollbackErrorLogging(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	orderUUID := uuid.New()

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM ` + postgres.OrdersTable).
		WithArgs(orderUUID.String()).
		WillReturnError(errors.New("exec error"))
	rollbackErr := errors.New("rollback failed")
	mock.ExpectRollback().WillReturnError(rollbackErr)

	// Capture log output
	var logOutput string
	log.SetOutput(&testLogWriter{&logOutput})
	defer log.SetOutput(nil)

	// Test
	err = repo.DeleteOrder(context.Background(), orderUUID)

	// Verify
	require.Error(t, err, "DeleteOrder should return an error")
	require.Contains(t, logOutput, "failed to rollback transaction", "Rollback error should be logged")
	require.Contains(t, logOutput, rollbackErr.Error(), "Rollback error message should be in logs")
	require.NoError(t, mock.ExpectationsWereMet())
}
