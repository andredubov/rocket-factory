package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/platform/pkg/logger"
)

func (r *ordersRepository) DeleteOrder(ctx context.Context, orderUUID uuid.UUID) error {
	return WithTx(ctx, r.pool, func(tx pgx.Tx) error {
		deleteBuilder := sq.Delete(OrdersTable).
			PlaceholderFormat(sq.Dollar).
			Where(sq.Eq{UUIDTableColumn: orderUUID})

		query, args, err := deleteBuilder.ToSql()
		if err != nil {
			logger.Error(ctx, "failed to build delete order query", zap.Error(err))
			return err
		}

		if _, err := tx.Exec(ctx, query, args...); err != nil {
			logger.Error(ctx, "failed to exec delete order query", zap.Error(err))
			return err
		}

		return nil
	})
}
