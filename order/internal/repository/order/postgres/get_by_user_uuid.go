package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

// GetUserOrders возвращает все заказы указанного пользователя вместе с их составными частями.
func (r *ordersRepository) GetUserOrders(ctx context.Context, userUUID uuid.UUID) ([]model.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			if err := tx.Rollback(ctx); err != nil {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}
	}()

	// Get orders
	orders, orderUUIDs, err := r.getUserOrderDetails(ctx, tx, userUUID)
	if err != nil {
		return nil, err
	}

	// Get parts if orders exist
	if len(orderUUIDs) > 0 {
		partsMap, err := r.getOrderPartsMap(ctx, tx, orderUUIDs)
		if err != nil {
			return nil, err
		}

		// Assign parts to orders
		for i := range orders {
			if parts, exists := partsMap[orders[i].OrderUUID]; exists {
				orders[i].PartUUIDs = parts
			}
		}
	}

	committed = true
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return orders, nil
}

// getUserOrderDetails получает основную информацию о заказах пользователя.
func (r *ordersRepository) getUserOrderDetails(ctx context.Context, tx pgx.Tx, userUUID uuid.UUID) ([]model.Order, []uuid.UUID, error) {
	query, args, err := sq.Select(
		uuidTableColumn,
		userUUIDTableColumn,
		totalPriceTableColumn,
		transactionUUIDTableColumn,
		paymentMethodTableColumn,
		statusTableColumn,
		createdAtTableColumn,
		updatedAtTableColumn,
	).
		From(ordersTable).
		Where(sq.Eq{userUUIDTableColumn: userUUID}).
		PlaceholderFormat(sq.Dollar).
		OrderBy(createdAtTableColumn + " DESC").
		ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build orders query: %w", err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query user orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	var orderUUIDs []uuid.UUID

	for rows.Next() {
		order, err := r.scanOrderRow(rows)
		if err != nil {
			return nil, nil, err
		}
		orders = append(orders, *order)
		orderUUIDs = append(orderUUIDs, order.OrderUUID)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error during orders iteration: %w", err)
	}

	return orders, orderUUIDs, nil
}

// scanOrderRow преобразует строку из БД в структуру Order.
func (r *ordersRepository) scanOrderRow(rows pgx.Rows) (*model.Order, error) {
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

	if err := rows.Scan(
		&dbOrderUUID,
		&dbUserUUID,
		&dbTotalPrice,
		&dbTransactionUUID,
		&dbPaymentMethod,
		&dbStatus,
		&dbCreatedAt,
		&dbUpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to scan order: %w", err)
	}

	orderStatus, err := model.NewOrderStatus(dbStatus)
	if err != nil {
		return nil, fmt.Errorf("invalid order status: %w", err)
	}

	paymentInfo := r.createPaymentInfo(dbTransactionUUID, dbPaymentMethod)

	return &model.Order{
		OrderUUID:   dbOrderUUID,
		UserUUID:    dbUserUUID,
		TotalPrice:  dbTotalPrice,
		PaymentInfo: paymentInfo,
		Status:      orderStatus,
		PartUUIDs:   []uuid.UUID{},
	}, nil
}

// createPaymentInfo создает структуру PaymentInfo на основе данных из БД.
func (r *ordersRepository) createPaymentInfo(transactionUUID *uuid.UUID, paymentMethod *string) *model.PaymentInfo {
	if transactionUUID == nil && paymentMethod == nil {
		return nil
	}

	info := &model.PaymentInfo{
		TransactionUUID: *transactionUUID,
	}

	if paymentMethod != nil {
		pm, err := model.NewPaymentMethod(*paymentMethod)
		if err != nil {
			return nil
		}
		info.PaymentMethod = pm
	}

	return info
}

// getOrderPartsMap получает все части для списка заказов одним запросом.
func (r *ordersRepository) getOrderPartsMap(ctx context.Context, tx pgx.Tx, orderUUIDs []uuid.UUID) (map[uuid.UUID][]uuid.UUID, error) {
	partsQuery, partsArgs, err := sq.Select(
		orderUUIDTableColumn,
		partUUIDTableColumn,
	).
		From(orderPartsTable).
		Where(sq.Eq{orderUUIDTableColumn: orderUUIDs}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build parts query: %w", err)
	}

	partsRows, err := tx.Query(ctx, partsQuery, partsArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to query order parts: %w", err)
	}
	defer partsRows.Close()

	partsMap := make(map[uuid.UUID][]uuid.UUID)
	for partsRows.Next() {
		var (
			orderUUID uuid.UUID
			partUUID  uuid.UUID
		)
		if err := partsRows.Scan(&orderUUID, &partUUID); err != nil {
			return nil, fmt.Errorf("failed to scan part: %w", err)
		}
		partsMap[orderUUID] = append(partsMap[orderUUID], partUUID)
	}

	if err := partsRows.Err(); err != nil {
		return nil, fmt.Errorf("error during parts iteration: %w", err)
	}

	return partsMap, nil
}
