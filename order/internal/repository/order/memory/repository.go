package memory

import (
	"sync"

	"github.com/google/uuid"

	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
	"github.com/andredubov/rocket-factory/order/internal/service"
)

// ordersRepository is an in-memory implementation of the Orders repository.
// It uses a sync.RWMutex for concurrent access protection and a map for storage.
type ordersRepository struct {
	mu     sync.RWMutex                   // Guards access to the orders map
	orders map[uuid.UUID]*repoModel.Order // Map storing orders by their UUID
}

// NewOrderRepository creates a new instance of an in-memory order repository.
// Returns an implementation of the repository.Orders interface.
func NewOrderRepository() service.OrdersRepository {
	return &ordersRepository{
		orders: make(map[uuid.UUID]*repoModel.Order), // Initialize empty orders map
	}
}
