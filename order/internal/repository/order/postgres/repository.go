package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/andredubov/rocket-factory/order/internal/service"
)

const (
	OrderPartsTable      = "order_parts"
	OrderUUIDTableColumn = "order_uuid"
	PartUUIDTableColumn  = "part_uuid"

	OrdersTable                = "orders"
	UUIDTableColumn            = "uuid"
	UserUUIDTableColumn        = "user_uuid"
	TotalPriceTableColumn      = "total_price"
	TransactionUUIDTableColumn = "transaction_uuid"
	PaymentMethodTableColumn   = "payment_method"
	StatusTableColumn          = "status"
	CreatedAtTableColumn       = "created_at"
	UpdatedAtTableColumn       = "updated_at"
)

// PgxPool
type PgxPool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

// ordersRepository is postgres implementation of the Orders repository.
type ordersRepository struct {
	pool PgxPool
}

// NewOrderRepository creates a new instance of an postgres order repository.
func NewOrderRepository(pool PgxPool) service.OrdersRepository {
	return &ordersRepository{
		pool: pool,
	}
}

// WithTx ...
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
