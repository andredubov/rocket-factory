package tests

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/order/postgres"
)

type testLogWriter struct {
	output *string
}

func (w *testLogWriter) Write(p []byte) (n int, err error) {
	*w.output += string(p)
	return len(p), nil
}

func TestAddOrder_Success(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)

	repoOrder := repoModel.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice: 100.50,
		PaymentInfo: &repoModel.PaymentInfo{
			TransactionUUID: uuid.New(),
			PaymentMethod:   repoModel.PaymentMethodCard,
		},
		Status: repoModel.OrderStatusPending,
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO `+postgres.OrdersTable).
		WithArgs(
			repoOrder.UserUUID,
			repoOrder.TotalPrice,
			repoOrder.PaymentInfo.TransactionUUID,
			repoOrder.PaymentInfo.PaymentMethod,
			repoOrder.Status,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// Expect inserts for all parts
	for _, partUUID := range repoOrder.PartUUIDs {
		mock.ExpectExec(`INSERT INTO `+postgres.OrderPartsTable).
			WithArgs(repoOrder.OrderUUID, partUUID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))
	}

	mock.ExpectCommit()

	// Test
	order := converter.OrderToModel(repoOrder) // метод покрыт юнит-тестами
	err = repo.AddOrder(context.Background(), *order)

	// Verify
	require.NoError(t, err, "AddOrder should not return an error")
	require.NoError(t, mock.ExpectationsWereMet(), "All expected SQL queries should be executed")
}

func TestAddOrder_EmptyParts(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)

	repoOrder := repoModel.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{},
		TotalPrice: 100.50,
		PaymentInfo: &repoModel.PaymentInfo{
			TransactionUUID: uuid.New(),
			PaymentMethod:   repoModel.PaymentMethodCard,
		},
		Status: repoModel.OrderStatusPending,
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO `+postgres.OrdersTable).
		WithArgs(
			repoOrder.UserUUID,
			repoOrder.TotalPrice,
			repoOrder.PaymentInfo.TransactionUUID,
			repoOrder.PaymentInfo.PaymentMethod,
			repoOrder.Status,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	// Test
	order := converter.OrderToModel(repoOrder)
	err = repo.AddOrder(context.Background(), *order)

	// Verify
	require.NoError(t, err, "AddOrder should not return an error for empty parts")
	require.NoError(t, mock.ExpectationsWereMet(), "All expected SQL queries should be executed")
}

func TestAddOrder_TransactionBeginError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)

	repoOrder := repoModel.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New()},
		TotalPrice: 100.50,
		PaymentInfo: &repoModel.PaymentInfo{
			TransactionUUID: uuid.New(),
			PaymentMethod:   repoModel.PaymentMethodCard,
		},
		Status: repoModel.OrderStatusPending,
	}

	// Expectations
	mock.ExpectBegin().WillReturnError(errors.New("begin error"))

	// Test
	order := converter.OrderToModel(repoOrder)
	err = repo.AddOrder(context.Background(), *order)

	// Verify
	require.Error(t, err, "AddOrder should return an error when transaction begin fails")
	require.Contains(t, err.Error(), "failed to begin transaction")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddOrder_OrderInsertError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)

	repoOrder := repoModel.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New()},
		TotalPrice: 100.50,
		PaymentInfo: &repoModel.PaymentInfo{
			TransactionUUID: uuid.New(),
			PaymentMethod:   repoModel.PaymentMethodCard,
		},
		Status: repoModel.OrderStatusPending,
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO `+postgres.OrdersTable).
		WithArgs(
			repoOrder.UserUUID,
			repoOrder.TotalPrice,
			repoOrder.PaymentInfo.TransactionUUID,
			repoOrder.PaymentInfo.PaymentMethod,
			repoOrder.Status,
		).
		WillReturnError(errors.New("insert error"))
	mock.ExpectRollback()

	// Test
	order := converter.OrderToModel(repoOrder)
	err = repo.AddOrder(context.Background(), *order)

	// Verify
	require.Error(t, err, "AddOrder should return an error when order insert fails")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddOrder_PartInsertError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)

	repoOrder := repoModel.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice: 100.50,
		PaymentInfo: &repoModel.PaymentInfo{
			TransactionUUID: uuid.New(),
			PaymentMethod:   repoModel.PaymentMethodCard,
		},
		Status: repoModel.OrderStatusPending,
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO `+postgres.OrdersTable).
		WithArgs(
			repoOrder.UserUUID,
			repoOrder.TotalPrice,
			repoOrder.PaymentInfo.TransactionUUID,
			repoOrder.PaymentInfo.PaymentMethod,
			repoOrder.Status,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// First part insert succeeds, second fails
	mock.ExpectExec(`INSERT INTO `+postgres.OrderPartsTable).
		WithArgs(repoOrder.OrderUUID, repoOrder.PartUUIDs[0]).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec(`INSERT INTO `+postgres.OrderPartsTable).
		WithArgs(repoOrder.OrderUUID, repoOrder.PartUUIDs[1]).
		WillReturnError(errors.New("part insert error"))
	mock.ExpectRollback()

	// Test
	order := converter.OrderToModel(repoOrder)
	err = repo.AddOrder(context.Background(), *order)

	// Verify
	require.Error(t, err, "AddOrder should return an error when part insert fails")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddOrder_CommitError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)

	repoOrder := repoModel.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New()},
		TotalPrice: 100.50,
		PaymentInfo: &repoModel.PaymentInfo{
			TransactionUUID: uuid.New(),
			PaymentMethod:   repoModel.PaymentMethodCard,
		},
		Status: repoModel.OrderStatusPending,
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO `+postgres.OrdersTable).
		WithArgs(
			repoOrder.UserUUID,
			repoOrder.TotalPrice,
			repoOrder.PaymentInfo.TransactionUUID,
			repoOrder.PaymentInfo.PaymentMethod,
			repoOrder.Status,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec(`INSERT INTO `+postgres.OrderPartsTable).
		WithArgs(repoOrder.OrderUUID, repoOrder.PartUUIDs[0]).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit().WillReturnError(errors.New("commit error"))

	// Test
	order := converter.OrderToModel(repoOrder)
	err = repo.AddOrder(context.Background(), *order)

	// Verify
	require.Error(t, err, "AddOrder should return an error when commit fails")
	require.Contains(t, err.Error(), "failed to commit transaction")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddOrder_RollbackErrorLogging(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)

	repoOrder := repoModel.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New()},
		TotalPrice: 100.50,
		PaymentInfo: &repoModel.PaymentInfo{
			TransactionUUID: uuid.New(),
			PaymentMethod:   repoModel.PaymentMethodCard,
		},
		Status: repoModel.OrderStatusPending,
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO `+postgres.OrdersTable).
		WithArgs(
			repoOrder.UserUUID,
			repoOrder.TotalPrice,
			repoOrder.PaymentInfo.TransactionUUID,
			repoOrder.PaymentInfo.PaymentMethod,
			repoOrder.Status,
		).
		WillReturnError(errors.New("order insert error"))

	// Force rollback to return an error
	rollbackErr := errors.New("rollback failed")
	mock.ExpectRollback().WillReturnError(rollbackErr)

	// Capture log output
	var logOutput string
	log.SetOutput(&testLogWriter{&logOutput})
	defer log.SetOutput(nil)

	// Test
	order := converter.OrderToModel(repoOrder)
	err = repo.AddOrder(context.Background(), *order)

	// Verify
	require.Error(t, err, "AddOrder should return an error")
	require.Contains(t, logOutput, "failed to rollback transaction", "Rollback error should be logged")
	require.Contains(t, logOutput, rollbackErr.Error(), "Rollback error message should be in logs")
	require.NoError(t, mock.ExpectationsWereMet())
}
