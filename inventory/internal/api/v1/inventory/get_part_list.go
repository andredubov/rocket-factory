package api

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/converter"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// ListParts handles requests to retrieve multiple inventory parts with optional filtering.
func (i *InventoryImplementation) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	// Convert gRPC filter to domain filter
	filter := converter.PartFilterFromListRequest(req)

	// Fetch parts from service using filter
	parts, err := i.inventoryService.GetPartList(ctx, filter)
	if err != nil {
		// if errors.Is(err, model.ErrPartNotFound) {
		// 	log.Printf("target parts not found")
		// 	return nil, status.Errorf(codes.NotFound, "target parts not found")
		// }
		return nil, err
	}

	// Convert domain models to gRPC response
	return converter.PartsToResponse(parts), nil
}
