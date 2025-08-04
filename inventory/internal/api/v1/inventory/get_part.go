package api

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/converter"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// GetPart handles requests to retrieve a single inventory part by its UUID.
func (i *InventoryImplementation) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	uuid := req.GetUuid() // Extract UUID from request

	// Fetch part from service
	part, err := i.inventoryService.GetPart(ctx, uuid)
	if err != nil {
		// if errors.Is(err, model.ErrPartNotFound) {
		// 	log.Printf("part with UUID %s not found", uuid)
		// 	return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", uuid)
		// }
		return nil, err
	}

	// Convert domain model to gRPC response
	return converter.PartToResponse(part), nil
}
