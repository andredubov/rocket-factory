package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/andredubov/rocket-factory/order/internal/model"
	"github.com/andredubov/rocket-factory/order/internal/repository"
)

func (r *ordersRepository) UpdateOrder(ctx context.Context, order model.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			if err := tx.Rollback(ctx); err != nil {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}
	}()

	// 1. Обновляем основную информацию о заказе
	updateBuilder := sq.Update(ordersTable).
		PlaceholderFormat(sq.Dollar).
		Set(totalPriceTableColumn, order.TotalPrice).
		Set(statusTableColumn, string(order.Status)).
		Where(sq.Eq{uuidTableColumn: order.OrderUUID})

	// Добавляем обновление платежной информации, если она есть
	if order.PaymentInfo != nil {
		updateBuilder = updateBuilder.
			Set(transactionUUIDTableColumn, order.PaymentInfo.TransactionUUID).
			Set(paymentMethodTableColumn, string(order.PaymentInfo.PaymentMethod))
	}

	updateBuilder = updateBuilder.Set(updatedAtTableColumn, time.Now())

	updateQuery, updateArgs, err := updateBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := tx.Exec(ctx, updateQuery, updateArgs...)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return repository.ErrOrderNotFoundWith(order.OrderUUID)
	}

	// 2. Обновляем состав заказа (части)
	// Сначала удаляем все существующие части заказа
	deletePartsQuery, deletePartsArgs, err := sq.Delete(orderPartsTable).
		Where(sq.Eq{orderUUIDTableColumn: order.OrderUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete parts query: %w", err)
	}

	if _, err := tx.Exec(ctx, deletePartsQuery, deletePartsArgs...); err != nil {
		return fmt.Errorf("failed to delete order parts: %w", err)
	}

	// Затем добавляем новые части заказа
	if len(order.PartUUIDs) > 0 {
		insertPartsBuilder := sq.Insert(orderPartsTable).
			Columns(orderUUIDTableColumn, partUUIDTableColumn).
			PlaceholderFormat(sq.Dollar)

		for _, partUUID := range order.PartUUIDs {
			insertPartsBuilder = insertPartsBuilder.Values(order.OrderUUID, partUUID)
		}

		insertPartsQuery, insertPartsArgs, err := insertPartsBuilder.ToSql()
		if err != nil {
			return fmt.Errorf("failed to build insert parts query: %w", err)
		}

		if _, err := tx.Exec(ctx, insertPartsQuery, insertPartsArgs...); err != nil {
			return fmt.Errorf("failed to insert order parts: %w", err)
		}
	}

	committed = true
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
