package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andredubov/rocket-factory/order/internal/service"
)

const (
	orderPartsTable      = "order_parts"
	orderUUIDTableColumn = "order_uuid"
	partUUIDTableColumn  = "part_uuid"

	ordersTable                = "orders"
	uuidTableColumn            = "uuid"
	userUUIDTableColumn        = "user_uuid"
	totalPriceTableColumn      = "total_price"
	transactionUUIDTableColumn = "transaction_uuid"
	paymentMethodTableColumn   = "payment_method"
	statusTableColumn          = "status"
	createdAtTableColumn       = "created_at"
	updatedAtTableColumn       = "updated_at"
)

// ordersRepository is postgres implementation of the Orders repository.
type ordersRepository struct {
	pool *pgxpool.Pool
}

// NewOrderRepository creates a new instance of an postgres order repository.
// Returns an implementation of the repository.Orders interface.
func NewOrderRepository(pool *pgxpool.Pool) service.OrdersRepository {
	return &ordersRepository{
		pool: pool,
	}
}
