package converter

import (
	"time"

	"github.com/andredubov/rocket-factory/order/internal/model"
	inventory_v1 "github.com/andredubov/rocket-factory/shared/pkg/proto/inventory/v1"
)

// PartToGetPartRequest converts a Part model to a GetPartRequest for gRPC calls to the inventory service.
func PartToGetPartRequest(part model.Part) *inventory_v1.GetPartRequest {
	return &inventory_v1.GetPartRequest{
		Uuid: part.Uuid,
	}
}

// PartFilterToListPartsRequest converts a PartFilter to a ListPartsRequest for gRPC calls.
func PartFilterToListPartsRequest(filter model.PartFilter) *inventory_v1.ListPartsRequest {
	partsFilter := &inventory_v1.PartsFilter{
		Uuids:                 filter.UUIDs,
		Names:                 filter.Names,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}

	// Конвертация категорий из model.Category в inventory_v1.Category
	if len(filter.Categories) > 0 {
		categories := make([]inventory_v1.Category, 0, len(filter.Categories))
		for _, cat := range filter.Categories {
			categories = append(categories, inventory_v1.Category(cat))
		}
		partsFilter.Categories = categories
	}

	return &inventory_v1.ListPartsRequest{
		Filter: partsFilter,
	}
}

// PartsFromListPartsResponse converts an inventory service response (ListPartsResponse) to a slice of Part models.
func PartsFromListPartsResponse(response *inventory_v1.ListPartsResponse) []model.Part {
	if response == nil || len(response.Parts) == 0 {
		return nil
	}

	parts := make([]model.Part, 0, len(response.Parts))
	for _, protoPart := range response.Parts {
		if protoPart == nil {
			continue
		}

		// Convert metadata
		metadata := make(map[string]model.Value)
		for k, v := range protoPart.Metadata {
			if v == nil {
				continue
			}
			val := model.Value{}
			switch x := v.Kind.(type) {
			case *inventory_v1.Value_StringValue:
				val.StringValue = &x.StringValue
			case *inventory_v1.Value_Int64Value:
				val.Int64Value = &x.Int64Value
			case *inventory_v1.Value_DoubleValue:
				val.DoubleValue = &x.DoubleValue
			case *inventory_v1.Value_BoolValue:
				val.BoolValue = &x.BoolValue
			}
			metadata[k] = val
		}

		// Convert timestamps
		var createdAt, updatedAt time.Time
		if protoPart.CreatedAt != nil {
			createdAt = protoPart.CreatedAt.AsTime()
		}
		if protoPart.UpdatedAt != nil {
			updatedAt = protoPart.UpdatedAt.AsTime()
		}

		part := model.Part{
			Uuid:          protoPart.Uuid,
			Name:          protoPart.Name,
			Description:   protoPart.Description,
			Price:         protoPart.Price,
			StockQuantity: protoPart.StockQuantity,
			Category:      model.PartCategory(protoPart.Category),
			Dimensions: model.Dimensions{
				Length: protoPart.Dimensions.GetLength(),
				Width:  protoPart.Dimensions.GetWidth(),
				Height: protoPart.Dimensions.GetHeight(),
				Weight: protoPart.Dimensions.GetWeight(),
			},
			Manufacturer: model.Manufacturer{
				Name:    protoPart.Manufacturer.GetName(),
				Country: protoPart.Manufacturer.GetCountry(),
				Website: protoPart.Manufacturer.GetWebsite(),
			},
			Tags:      protoPart.Tags,
			Metadata:  metadata,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		parts = append(parts, part)
	}

	return parts
}
