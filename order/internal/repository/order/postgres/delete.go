package postgres

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *ordersRepository) DeleteOrder(ctx context.Context, uuid uuid.UUID) error {
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

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true

	return nil
}
