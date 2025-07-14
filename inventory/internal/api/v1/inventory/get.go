package server

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/andredubov/rocket-factory/inventory/internal/converter"
	"github.com/andredubov/rocket-factory/inventory/internal/model"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// GetPart handles requests to retrieve a single inventory part by its UUID.
func (i *InventoryImplementation) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	uuid := req.GetUuid() // Extract UUID from request

	// Fetch part from service
	part, err := i.inventoryService.GetPart(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			log.Printf("part with UUID %s not found", uuid)
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", uuid)
		}
		return nil, err
	}

	// Convert domain model to gRPC response
	return converter.PartToResponse(part), nil
}

// ListParts handles requests to retrieve multiple inventory parts with optional filtering.
func (i *InventoryImplementation) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	// Convert gRPC filter to domain filter
	filter := converter.PartFilterFromListRequest(req)

	// Fetch parts from service using filter
	parts, err := i.inventoryService.GetPartList(ctx, filter)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			log.Printf("target parts not found")
			return nil, status.Errorf(codes.NotFound, "target parts not found")
		}
		return nil, err
	}

	// Convert domain models to gRPC response
	return converter.PartsToResponse(parts), nil
}
