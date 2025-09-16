package tests

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/converter"
	repoOrder "github.com/andredubov/rocket-factory/order/internal/repository/model"
	"github.com/andredubov/rocket-factory/order/internal/repository/order/postgres"
)

func createTestOrder() repoOrder.Order {
	return repoOrder.Order{
		OrderUUID:  uuid.New(),
		UserUUID:   uuid.New(),
		PartUUIDs:  []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice: 100.50,
		PaymentInfo: &repoOrder.PaymentInfo{
			TransactionUUID: uuid.New(),
			PaymentMethod:   repoOrder.PaymentMethodCard,
		},
		Status: repoOrder.OrderStatusPaid,
	}
}

func TestUpdateOrder_Success(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	repoOrder := createTestOrder()

	// Expectations
	mock.ExpectBegin()

	// Ожидаем обновление основной информации о заказе
	mock.ExpectExec(`UPDATE `+postgres.OrdersTable+` SET`).
		WithArgs(
			repoOrder.TotalPrice,                        // total_price (float64)
			string(repoOrder.Status),                    // status (string)
			repoOrder.PaymentInfo.TransactionUUID,       // transaction_uuid (string)
			string(repoOrder.PaymentInfo.PaymentMethod), // payment_method (string)
			pgxmock.AnyArg(),                            // updated_at (time.Time)
			repoOrder.OrderUUID.String(),                // where uuid (string)
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Ожидаем удаление старых частей заказа
	mock.ExpectExec(`DELETE FROM ` + postgres.OrderPartsTable).
		WithArgs(repoOrder.OrderUUID.String()). // order_uuid (string)
		WillReturnResult(pgxmock.NewResult("DELETE", 2))

	// Ожидаем добавление новых частей заказа
	mock.ExpectExec(`INSERT INTO `+postgres.OrderPartsTable).
		WithArgs(
			repoOrder.OrderUUID, repoOrder.PartUUIDs[0],
			repoOrder.OrderUUID, repoOrder.PartUUIDs[1],
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 2))

	mock.ExpectCommit()

	// Test
	order := converter.OrderToModel(repoOrder)

	updateInfo := model.OrderUpdateInfo{
		OrderUUID:   order.OrderUUID,
		UserUUID:    &order.UserUUID,
		PartUUIDs:   order.PartUUIDs,
		TotalPrice:  &order.TotalPrice,
		PaymentInfo: order.PaymentInfo,
		Status:      &order.Status,
	}

	err = repo.UpdateOrder(context.Background(), updateInfo)

	// Verify
	require.NoError(t, err, "UpdateOrder should not return an error")
	require.NoError(t, mock.ExpectationsWereMet(), "All expected SQL queries should be executed")
}

func TestUpdateOrder_WithoutPaymentInfo(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	repoOrder := createTestOrder()
	repoOrder.PaymentInfo = nil // Заказ без платежной информации

	// Expectations
	mock.ExpectBegin()

	// Ожидаем обновление основной информации о заказе
	mock.ExpectExec(`UPDATE `+postgres.OrdersTable+` SET`).
		WithArgs(
			repoOrder.TotalPrice,         // total_price (float64)
			string(repoOrder.Status),     // status (string)
			pgxmock.AnyArg(),             // updated_at (time.Time)
			repoOrder.OrderUUID.String(), // where uuid (string)
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Ожидаем удаление старых частей заказа
	mock.ExpectExec(`DELETE FROM order_parts`).
		WithArgs(repoOrder.OrderUUID.String()).
		WillReturnResult(pgxmock.NewResult("DELETE", 2))

	// Ожидаем добавление новых частей заказа
	mock.ExpectExec(`INSERT INTO order_parts`).
		WithArgs(
			repoOrder.OrderUUID, repoOrder.PartUUIDs[0],
			repoOrder.OrderUUID, repoOrder.PartUUIDs[1],
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 2))
	mock.ExpectCommit()

	// Test
	order := converter.OrderToModel(repoOrder)

	updateInfo := model.OrderUpdateInfo{
		OrderUUID:   order.OrderUUID,
		UserUUID:    &order.UserUUID,
		PartUUIDs:   order.PartUUIDs,
		TotalPrice:  &order.TotalPrice,
		PaymentInfo: order.PaymentInfo,
		Status:      &order.Status,
	}

	err = repo.UpdateOrder(context.Background(), updateInfo)

	// Verify
	require.NoError(t, err, "UpdateOrder should not return an error")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateOrder_NoParts(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	repoOrder := createTestOrder()
	repoOrder.PartUUIDs = []uuid.UUID{} // Заказ без частей

	// Expectations
	mock.ExpectBegin()

	// Ожидаем обновление основной информации о заказе
	mock.ExpectExec(`UPDATE `+postgres.OrdersTable+` SET`).
		WithArgs(
			repoOrder.TotalPrice,                        // total_price (float64)
			string(repoOrder.Status),                    // status (string)
			repoOrder.PaymentInfo.TransactionUUID,       // transaction_uuid (string)
			string(repoOrder.PaymentInfo.PaymentMethod), // payment_method (string)
			pgxmock.AnyArg(),                            // updated_at (time.Time)
			repoOrder.OrderUUID.String(),                // where uuid (string)
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Ожидаем удаление старых частей заказа
	mock.ExpectExec(`DELETE FROM order_parts`).
		WithArgs(repoOrder.OrderUUID.String()).
		WillReturnResult(pgxmock.NewResult("DELETE", 2))

	// Не добавляем новые части заказа

	mock.ExpectCommit()

	// Test
	order := converter.OrderToModel(repoOrder)

	updateInfo := model.OrderUpdateInfo{
		OrderUUID:   order.OrderUUID,
		UserUUID:    &order.UserUUID,
		PartUUIDs:   order.PartUUIDs,
		TotalPrice:  &order.TotalPrice,
		PaymentInfo: order.PaymentInfo,
		Status:      &order.Status,
	}

	err = repo.UpdateOrder(context.Background(), updateInfo)

	// Verify
	require.NoError(t, err, "UpdateOrder should not return an error")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateOrder_NotFound(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	repoOrder := createTestOrder()

	// Expectations
	mock.ExpectBegin()

	// Ожидаем обновление основной информации о заказе
	mock.ExpectExec(`UPDATE `+postgres.OrdersTable+` SET`).
		WithArgs(
			repoOrder.TotalPrice,                        // total_price (float64)
			string(repoOrder.Status),                    // status (string)
			repoOrder.PaymentInfo.TransactionUUID,       // transaction_uuid (string)
			string(repoOrder.PaymentInfo.PaymentMethod), // payment_method (string)
			pgxmock.AnyArg(),                            // updated_at (time.Time)
			repoOrder.OrderUUID.String(),                // where uuid (string)
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	mock.ExpectRollback()

	// Test
	order := converter.OrderToModel(repoOrder)

	updateInfo := model.OrderUpdateInfo{
		OrderUUID:   order.OrderUUID,
		UserUUID:    &order.UserUUID,
		PartUUIDs:   order.PartUUIDs,
		TotalPrice:  &order.TotalPrice,
		PaymentInfo: order.PaymentInfo,
		Status:      &order.Status,
	}

	err = repo.UpdateOrder(context.Background(), updateInfo)

	// Verify
	require.Error(t, err, "UpdateOrder should return an error")
	require.Equal(t, err, model.ErrOrderNotFoundWith(order.OrderUUID))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateOrder_BeginTxError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	testOrder := createTestOrder()

	// Expectations
	expectedErr := errors.New("begin error")
	mock.ExpectBegin().WillReturnError(expectedErr)

	// Test
	order := converter.OrderToModel(testOrder)

	updateInfo := model.OrderUpdateInfo{
		OrderUUID:   order.OrderUUID,
		UserUUID:    &order.UserUUID,
		PartUUIDs:   order.PartUUIDs,
		TotalPrice:  &order.TotalPrice,
		PaymentInfo: order.PaymentInfo,
		Status:      &order.Status,
	}

	err = repo.UpdateOrder(context.Background(), updateInfo)

	// Verify
	require.Error(t, err, "UpdateOrder should return an error")
	require.Contains(t, err.Error(), "failed to begin transaction")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateOrder_UpdateError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	repoOrder := createTestOrder()

	// Expectations
	mock.ExpectBegin()

	// Ожидаем обновление основной информации о заказе
	mock.ExpectExec(`UPDATE `+postgres.OrdersTable+` SET`).
		WithArgs(
			repoOrder.TotalPrice,                        // total_price (float64)
			string(repoOrder.Status),                    // status (string)
			repoOrder.PaymentInfo.TransactionUUID,       // transaction_uuid (string)
			string(repoOrder.PaymentInfo.PaymentMethod), // payment_method (string)
			pgxmock.AnyArg(),                            // updated_at (time.Time)
			repoOrder.OrderUUID.String(),                // where uuid (string)
		).
		WillReturnError(errors.New("update error"))

	mock.ExpectRollback()

	// Test
	order := converter.OrderToModel(repoOrder)

	updateInfo := model.OrderUpdateInfo{
		OrderUUID:   order.OrderUUID,
		UserUUID:    &order.UserUUID,
		PartUUIDs:   order.PartUUIDs,
		TotalPrice:  &order.TotalPrice,
		PaymentInfo: order.PaymentInfo,
		Status:      &order.Status,
	}

	err = repo.UpdateOrder(context.Background(), updateInfo)

	// Verify
	require.Error(t, err, "UpdateOrder should return an error")
	// require.Contains(t, err.Error(), "failed to update order")
	require.Contains(t, err.Error(), "failed to execute update order")
}

func TestUpdateOrder_DeletePartsError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	repoOrder := createTestOrder()

	// Expectations
	mock.ExpectBegin()

	// Ожидаем обновление основной информации о заказе
	mock.ExpectExec(`UPDATE `+postgres.OrdersTable+` SET`).
		WithArgs(
			repoOrder.TotalPrice,                        // total_price (float64)
			string(repoOrder.Status),                    // status (string)
			repoOrder.PaymentInfo.TransactionUUID,       // transaction_uuid (string)
			string(repoOrder.PaymentInfo.PaymentMethod), // payment_method (string)
			pgxmock.AnyArg(),                            // updated_at (time.Time)
			repoOrder.OrderUUID.String(),                // where uuid (string)
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Ожидаем удаление старых частей заказа
	mock.ExpectExec(`DELETE FROM ` + postgres.OrderPartsTable).
		WithArgs(repoOrder.OrderUUID.String()). // order_uuid (string)
		WillReturnError(errors.New("delete error"))

	mock.ExpectRollback()

	// Test
	order := converter.OrderToModel(repoOrder)

	updateInfo := model.OrderUpdateInfo{
		OrderUUID:   order.OrderUUID,
		UserUUID:    &order.UserUUID,
		PartUUIDs:   order.PartUUIDs,
		TotalPrice:  &order.TotalPrice,
		PaymentInfo: order.PaymentInfo,
		Status:      &order.Status,
	}

	err = repo.UpdateOrder(context.Background(), updateInfo)

	// Verify
	require.Error(t, err, "UpdateOrder should return an error")
	require.Contains(t, err.Error(), "failed to delete order parts")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateOrder_InsertPartsError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	// Capture log output
	var logOutput string
	log.SetOutput(&testLogWriter{&logOutput})
	defer log.SetOutput(nil)

	repo := postgres.NewOrderRepository(mock)
	repoOrder := createTestOrder()

	// Expectations
	mock.ExpectBegin()

	// Ожидаем обновление основной информации о заказе
	mock.ExpectExec(`UPDATE `+postgres.OrdersTable+` SET`).
		WithArgs(
			repoOrder.TotalPrice,                        // total_price (float64)
			string(repoOrder.Status),                    // status (string)
			repoOrder.PaymentInfo.TransactionUUID,       // transaction_uuid (string)
			string(repoOrder.PaymentInfo.PaymentMethod), // payment_method (string)
			pgxmock.AnyArg(),                            // updated_at (time.Time)
			repoOrder.OrderUUID.String(),                // where uuid (string)
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Ожидаем удаление старых частей заказа
	mock.ExpectExec(`DELETE FROM ` + postgres.OrderPartsTable).
		WithArgs(repoOrder.OrderUUID.String()). // order_uuid (string)
		WillReturnResult(pgxmock.NewResult("DELETE", int64(len(repoOrder.PartUUIDs))))

	// Ожидаем добавление новых частей заказа
	mock.ExpectExec(`INSERT INTO `+postgres.OrderPartsTable).
		WithArgs(
			repoOrder.OrderUUID, repoOrder.PartUUIDs[0],
			repoOrder.OrderUUID, repoOrder.PartUUIDs[1],
		).
		WillReturnError(errors.New("delete error"))

	// Test
	order := converter.OrderToModel(repoOrder)

	updateInfo := model.OrderUpdateInfo{
		OrderUUID:   order.OrderUUID,
		UserUUID:    &order.UserUUID,
		PartUUIDs:   order.PartUUIDs,
		TotalPrice:  &order.TotalPrice,
		PaymentInfo: order.PaymentInfo,
		Status:      &order.Status,
	}

	err = repo.UpdateOrder(context.Background(), updateInfo)

	// Verify
	require.Error(t, err, "UpdateOrder should return an error")
	require.Contains(t, err.Error(), "failed to insert order parts")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateOrder_CommitError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	repoOrder := createTestOrder()

	// Expectations
	mock.ExpectBegin()

	// Ожидаем обновление основной информации о заказе
	mock.ExpectExec(`UPDATE `+postgres.OrdersTable+` SET`).
		WithArgs(
			repoOrder.TotalPrice,                        // total_price (float64)
			string(repoOrder.Status),                    // status (string)
			repoOrder.PaymentInfo.TransactionUUID,       // transaction_uuid (string)
			string(repoOrder.PaymentInfo.PaymentMethod), // payment_method (string)
			pgxmock.AnyArg(),                            // updated_at (time.Time)
			repoOrder.OrderUUID.String(),                // where uuid (string)
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Ожидаем удаление старых частей заказа
	mock.ExpectExec(`DELETE FROM ` + postgres.OrderPartsTable).
		WithArgs(repoOrder.OrderUUID.String()). // order_uuid (string)
		WillReturnResult(pgxmock.NewResult("DELETE", 2))

	// Ожидаем добавление новых частей заказа
	mock.ExpectExec(`INSERT INTO `+postgres.OrderPartsTable).
		WithArgs(
			repoOrder.OrderUUID, repoOrder.PartUUIDs[0],
			repoOrder.OrderUUID, repoOrder.PartUUIDs[1],
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 2))

	mock.ExpectCommit().WillReturnError(errors.New("commit error"))

	// Capture log output
	var logOutput string
	log.SetOutput(&testLogWriter{&logOutput})
	defer log.SetOutput(nil)

	// Test
	order := converter.OrderToModel(repoOrder)

	updateInfo := model.OrderUpdateInfo{
		OrderUUID:   order.OrderUUID,
		UserUUID:    &order.UserUUID,
		PartUUIDs:   order.PartUUIDs,
		TotalPrice:  &order.TotalPrice,
		PaymentInfo: order.PaymentInfo,
		Status:      &order.Status,
	}

	err = repo.UpdateOrder(context.Background(), updateInfo)

	// Verify
	require.Error(t, err, "UpdateOrder should return an error")
	require.Contains(t, err.Error(), "failed to commit transaction")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateOrder_RollbackError(t *testing.T) {
	// Setup
	mock, err := pgxmock.NewPool()
	require.NoError(t, err, "Failed to create mock pool")
	defer mock.Close()

	repo := postgres.NewOrderRepository(mock)
	repoOrder := createTestOrder()

	// Expectations
	mock.ExpectBegin()

	// Ожидаем обновление основной информации о заказе
	mock.ExpectExec(`UPDATE `+postgres.OrdersTable+` SET`).
		WithArgs(
			repoOrder.TotalPrice,                        // total_price (float64)
			string(repoOrder.Status),                    // status (string)
			repoOrder.PaymentInfo.TransactionUUID,       // transaction_uuid (string)
			string(repoOrder.PaymentInfo.PaymentMethod), // payment_method (string)
			pgxmock.AnyArg(),                            // updated_at (time.Time)
			repoOrder.OrderUUID.String(),                // where uuid (string)
		).
		WillReturnError(errors.New("update error"))

	rollbackErr := errors.New("rollback error")
	mock.ExpectRollback().WillReturnError(rollbackErr)

	// Capture log output
	var logOutput string
	log.SetOutput(&testLogWriter{&logOutput})
	defer log.SetOutput(nil)

	// Test
	order := converter.OrderToModel(repoOrder)

	updateInfo := model.OrderUpdateInfo{
		OrderUUID:   order.OrderUUID,
		UserUUID:    &order.UserUUID,
		PartUUIDs:   order.PartUUIDs,
		TotalPrice:  &order.TotalPrice,
		PaymentInfo: order.PaymentInfo,
		Status:      &order.Status,
	}

	err = repo.UpdateOrder(context.Background(), updateInfo)

	// Verify
	require.Error(t, err, "UpdateOrder should return an error")
	require.Contains(t, logOutput, "failed to rollback transaction", "Rollback error should be logged")
	require.Contains(t, logOutput, rollbackErr.Error(), "Rollback error message should be in logs")
	require.NoError(t, mock.ExpectationsWereMet())
}
