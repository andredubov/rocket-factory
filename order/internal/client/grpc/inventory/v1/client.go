package inventory

import (
	"github.com/andredubov/rocket-factory/order/internal/client/grpc"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// inventoryClient is a client for interacting with the Inventory Service.
// It wraps the auto-generated gRPC client to provide inventory-specific functionality.
type inventoryClient struct {
	generatedClient inventory_v1.InventoryServiceClient
}

// NewClient creates a new instance of inventoryClient.
func NewClient(client inventory_v1.InventoryServiceClient) grpc.InventoryClient {
	return &inventoryClient{
		generatedClient: client,
	}
}
