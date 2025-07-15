package orders

import (
	"context"

	"github.com/google/uuid"
)

// DeleteOrder removes an order by its UUID.
func (s *ordersService) DeleteOrder(ctx context.Context, uuid uuid.UUID) error {
	return s.ordersRepository.DeleteOrder(ctx, uuid)
}
