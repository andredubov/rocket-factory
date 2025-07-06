package server

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/andredubov/rocket-factory/inventory/internal/repository/converter"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// GetPart handles requests to retrieve a single inventory part by its UUID.
// It converts the gRPC request to a domain model, fetches the part from the repository,
// and converts the result back to a gRPC response format.
// Returns:
// - gRPC NotFound error if the part doesn't exist
// - gRPC Internal error for repository failures
// - Successful response with the part data if found
func (i *Implementation) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	uuid := req.GetUuid() // Extract UUID from request

	// Fetch part from repository
	part, err := i.inventoryRepository.GetPart(ctx, uuid)
	if err != nil {
		log.Printf("target part not found: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "target part not found")
	}

	// Convert domain model to gRPC response
	return converter.PartToResponse(part), nil
}

// ListParts handles requests to retrieve multiple inventory parts with optional filtering.
// It converts the gRPC filter request to a domain filter, fetches matching parts
// from the repository, and converts the results to a gRPC response format.
// Returns:
// - gRPC Internal error for repository failures
// - Successful response with the list of parts (empty if no matches)
func (i *Implementation) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	// Convert gRPC filter to domain filter
	filter := converter.PartFilterFromListRequest(req)

	// Fetch parts from repository using filter
	parts, err := i.inventoryRepository.GetPartList(ctx, filter)
	if err != nil {
		log.Printf("target parts not found: %s", err.Error())
		return nil, status.Errorf(codes.Internal, "target parts not found")
	}

	// Convert domain models to gRPC response
	return converter.PartsToResponse(parts), nil
}
