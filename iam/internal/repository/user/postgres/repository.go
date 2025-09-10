package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/andredubov/rocket-factory/iam/internal/service"
)

const (
	tableUsers               = "users"
	tableNotificationMethods = "notification_methods"

	columnUUID         = "uuid"
	columnLogin        = "login"
	columnEmail        = "email"
	columnPassword     = "password"
	columnCreatedAt    = "created_at"
	columnUpdatedAt    = "updated_at"
	columnUserUUID     = "user_uuid"
	columnProviderName = "provider_name"
	columnTarget       = "target"
)

type PgxPool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type usersRepository struct {
	pool PgxPool
}

func NewUsersRepository(pool PgxPool) service.UsersRepository {
	return &usersRepository{
		pool: pool,
	}
}

func WithTx(ctx context.Context, pool PgxPool, action func(tx pgx.Tx) error) error {
	committed := false
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if !committed {
			if err := tx.Rollback(ctx); err != nil {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}
	}()

	if err = action(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true
	return nil
}
