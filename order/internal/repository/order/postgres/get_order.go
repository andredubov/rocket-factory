package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository"
)

// GetOrder возвращает заказ по его UUID вместе со всеми связанными данными.
// Работает в рамках транзакции для обеспечения целостности данных.
func (r *ordersRepository) GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction:: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			if err := tx.Rollback(ctx); err != nil {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}
	}()

	// Получаем основную информацию о заказе
	order, err := r.getOrderDetails(ctx, tx, orderUUID)
	if err != nil {
		return nil, err
	}

	// Получаем части заказа
	order.PartUUIDs, err = r.getOrderParts(ctx, tx, order.OrderUUID)
	if err != nil {
		return nil, err
	}

	committed = true
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

// getOrderDetails получает основную информацию о заказе из БД
func (r *ordersRepository) getOrderDetails(ctx context.Context, tx pgx.Tx, orderUUID uuid.UUID) (*model.Order, error) {
	query, args, err := sq.Select(
		UUIDTableColumn,
		UserUUIDTableColumn,
		TotalPriceTableColumn,
		TransactionUUIDTableColumn,
		PaymentMethodTableColumn,
		StatusTableColumn,
		CreatedAtTableColumn,
		UpdatedAtTableColumn,
	).
		From(OrdersTable).
		Where(sq.Eq{UUIDTableColumn: orderUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var (
		dbOrderUUID       uuid.UUID
		dbUserUUID        uuid.UUID
		dbTotalPrice      float64
		dbTransactionUUID *uuid.UUID
		dbPaymentMethod   *string
		dbStatus          string
		dbCreatedAt       time.Time
		dbUpdatedAt       *time.Time
	)

	err = tx.QueryRow(ctx, query, args...).Scan(
		&dbOrderUUID,
		&dbUserUUID,
		&dbTotalPrice,
		&dbTransactionUUID,
		&dbPaymentMethod,
		&dbStatus,
		&dbCreatedAt,
		&dbUpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrOrderNotFoundWith(orderUUID)
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return r.buildOrderModel(dbOrderUUID, dbUserUUID, dbTotalPrice, dbTransactionUUID, dbPaymentMethod, dbStatus)
}

// getOrderParts получает список UUID частей для указанного заказа
func (r *ordersRepository) getOrderParts(ctx context.Context, tx pgx.Tx, orderUUID uuid.UUID) ([]uuid.UUID, error) {
	query, args, err := sq.Select(PartUUIDTableColumn).
		From(OrderPartsTable).
		Where(sq.Eq{OrderUUIDTableColumn: orderUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build parts query: %w", err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute parts query: %w", err)
	}
	defer rows.Close()

	var parts []uuid.UUID
	for rows.Next() {
		var partUUID uuid.UUID
		if err := rows.Scan(&partUUID); err != nil {
			return nil, fmt.Errorf("failed to scan part UUID: %w", err)
		}
		parts = append(parts, partUUID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing results: %w", err)
	}

	return parts, nil
}

// buildOrderModel создает объект Order на основе данных из БД
func (r *ordersRepository) buildOrderModel(
	orderUUID uuid.UUID,
	userUUID uuid.UUID,
	totalPrice float64,
	transactionUUID *uuid.UUID,
	paymentMethod *string,
	status string,
) (*model.Order, error) {
	orderStatus, err := model.NewOrderStatus(status)
	if err != nil {
		return nil, fmt.Errorf("invalid order status: %w", err)
	}

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		TotalPrice: totalPrice,
		Status:     orderStatus,
		PartUUIDs:  []uuid.UUID{},
	}

	if transactionUUID != nil || paymentMethod != nil {
		order.PaymentInfo = &model.PaymentInfo{
			TransactionUUID: *transactionUUID,
		}

		if paymentMethod != nil {
			pm, err := model.NewPaymentMethod(*paymentMethod)
			if err != nil {
				return nil, fmt.Errorf("invalid payment method: %w", err)
			}
			order.PaymentInfo.PaymentMethod = pm
		}
	}

	return order, nil
}
