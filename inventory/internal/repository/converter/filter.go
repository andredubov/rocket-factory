package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// PartFilterFromListRequest converts a gRPC ListPartsRequest filter to a domain PartFilter
// Handles nil checks and converts protobuf Categories to domain Categories
func PartFilterFromListRequest(r *inventory_v1.ListPartsRequest) model.PartFilter {
	// Return empty filter if request or filter is nil
	if r == nil || r.Filter == nil {
		return model.PartFilter{}
	}

	// Convert protobuf categories to domain categories
	pbCategories := r.GetFilter().GetCategories()
	categories := make([]model.Category, len(pbCategories))
	for i, cat := range pbCategories {
		categories[i] = model.Category{
			ID:   cat.GetId(),
			Name: cat.GetName(),
		}
	}

	// Build and return the domain filter with all converted fields
	return model.PartFilter{
		UUIDs:                 r.GetFilter().GetUuids(),
		Names:                 r.GetFilter().GetNames(),
		Categories:            categories,
		ManufacturerCountries: r.GetFilter().GetManufacturerCountries(),
		Tags:                  r.GetFilter().GetTags(),
	}
}

// PartsToResponse converts a slice of domain Parts to a gRPC ListPartsResponse
// Handles conversion of all part fields including nested structures and timestamps
func PartsToResponse(parts []model.Part) *inventory_v1.ListPartsResponse {
	// Pre-allocate slice for protobuf parts
	pbParts := make([]*inventory_v1.Part, len(parts))

	// Convert each domain Part to protobuf Part
	for i, part := range parts {
		pbParts[i] = &inventory_v1.Part{
			Uuid:          part.Uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Category: &inventory_v1.Category{
				Id:   part.Category.ID,
				Name: part.Category.Name,
			},
			Dimensions: &inventory_v1.Dimensions{
				Length: part.Dimensions.Length,
				Width:  part.Dimensions.Width,
				Height: part.Dimensions.Height,
				Weight: part.Dimensions.Weight,
			},
			Manufacturer: &inventory_v1.Manufacturer{
				Name:    part.Manufacturer.Name,
				Country: part.Manufacturer.Country,
				Website: part.Manufacturer.Website,
			},
			Tags:      part.Tags,
			CreatedAt: timestamppb.New(part.CreatedAt), // Convert time.Time to Timestamp
			UpdatedAt: timestamppb.New(part.UpdatedAt), // Convert time.Time to Timestamp
		}
	}

	return &inventory_v1.ListPartsResponse{
		Parts: pbParts,
	}
}
