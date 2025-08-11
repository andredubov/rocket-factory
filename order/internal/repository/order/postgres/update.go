package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/andredubov/rocket-factory/order/internal/model"
)

func (r *ordersRepository) UpdateOrder(ctx context.Context, order model.Order) error {
	return WithTx(ctx, r.pool, func(tx pgx.Tx) error {
		// 1. Обновляем основную информацию о заказе
		updateBuilder := sq.Update(OrdersTable).
			PlaceholderFormat(sq.Dollar).
			Set(TotalPriceTableColumn, order.TotalPrice).
			Set(StatusTableColumn, string(order.Status)).
			Where(sq.Eq{UUIDTableColumn: order.OrderUUID})

		// Добавляем обновление платежной информации, если она есть
		if order.PaymentInfo != nil {
			updateBuilder = updateBuilder.
				Set(TransactionUUIDTableColumn, order.PaymentInfo.TransactionUUID).
				Set(PaymentMethodTableColumn, string(order.PaymentInfo.PaymentMethod))
		}

		updateBuilder = updateBuilder.Set(UpdatedAtTableColumn, time.Now())

		updateQuery, updateArgs, err := updateBuilder.ToSql()
		if err != nil {
			return fmt.Errorf("failed to build update query: %w", err)
		}

		result, err := tx.Exec(ctx, updateQuery, updateArgs...)
		if err != nil {
			return fmt.Errorf("failed to update order: %w", err)
		}

		if result.RowsAffected() == 0 {
			return model.ErrOrderNotFoundWith(order.OrderUUID)
		}

		// 2. Обновляем состав заказа (части)
		// Сначала удаляем все существующие части заказа
		deletePartsQuery, deletePartsArgs, err := sq.Delete(OrderPartsTable).
			Where(sq.Eq{OrderUUIDTableColumn: order.OrderUUID}).
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
			insertPartsBuilder := sq.Insert(OrderPartsTable).
				Columns(OrderUUIDTableColumn, PartUUIDTableColumn).
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
		return nil
	})
}
