package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// PartFromRequest converts a gRPC GetPartRequest to a domain Part
func PartFromRequest(request *inventory_v1.GetPartRequest) model.Part {
	return model.Part{
		Uuid: request.GetUuid(),
	}
}

// PartToResponse converts a domain Part to a gRPC GetPartResponse
func PartToResponse(part *model.Part) *inventory_v1.GetPartResponse {
	// Convert metadata map from domain Value to protobuf Value
	metadata := make(map[string]*inventory_v1.Value)
	for k, v := range part.Metadata {
		value := &inventory_v1.Value{}

		// Handle each possible value type in the oneof
		switch {
		case v.StringValue != nil:
			value.Kind = &inventory_v1.Value_StringValue{StringValue: *v.StringValue}
		case v.Int64Value != nil:
			value.Kind = &inventory_v1.Value_Int64Value{Int64Value: *v.Int64Value}
		case v.DoubleValue != nil:
			value.Kind = &inventory_v1.Value_DoubleValue{DoubleValue: *v.DoubleValue}
		case v.BoolValue != nil:
			value.Kind = &inventory_v1.Value_BoolValue{BoolValue: *v.BoolValue}
		}

		metadata[k] = value
	}

	// Build and return the complete response with all converted fields
	return &inventory_v1.GetPartResponse{
		Part: &inventory_v1.Part{
			Uuid:          part.Uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Dimensions: &inventory_v1.Dimensions{
				Length: part.Dimensions.Length,
				Width:  part.Dimensions.Width,
				Height: part.Dimensions.Height,
				Weight: part.Dimensions.Weight,
			},
			Category: inventory_v1.Category(part.Category),
			Manufacturer: &inventory_v1.Manufacturer{
				Name:    part.Manufacturer.Name,
				Country: part.Manufacturer.Country,
				Website: part.Manufacturer.Website,
			},
			Metadata:  metadata,
			Tags:      part.Tags,
			CreatedAt: timestamppb.New(part.CreatedAt), // Convert time.Time to Timestamp
			UpdatedAt: timestamppb.New(part.UpdatedAt), // Convert time.Time to Timestamp
		},
	}
}
