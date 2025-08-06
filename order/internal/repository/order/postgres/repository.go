package postgres

import (
	"context"

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
