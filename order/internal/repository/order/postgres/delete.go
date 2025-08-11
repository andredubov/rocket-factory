package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *ordersRepository) DeleteOrder(ctx context.Context, uuid uuid.UUID) error {
	return WithTx(ctx, r.pool, func(tx pgx.Tx) error {
		deleteBuilder := sq.Delete(OrdersTable).
			PlaceholderFormat(sq.Dollar).
			Where(sq.Eq{OrderUUIDTableColumn: uuid})

		query, args, err := deleteBuilder.ToSql()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, query, args...); err != nil {
			return err
		}

		return nil
	})
}
