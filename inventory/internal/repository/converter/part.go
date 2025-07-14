package converter

import (
	"github.com/andredubov/rocket-factory/inventory/internal/model"
	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// PartToRepoModel converts a domain model (repo.Part) to a repository model (repoModel.Part).
func PartToRepoModel(part model.Part) repoModel.Part {
	return repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      repoModel.PartCategory(part.Category),
		Dimensions: repoModel.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      part.Tags,
		Metadata:  convertMetadataToRepo(part.Metadata),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

// PartToModel converts a repository model (repoModel.Part) to a domain model (model.Part).
func PartToModel(part repoModel.Part) model.Part {
	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.PartCategory(part.Category),
		Dimensions: model.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      part.Tags,
		Metadata:  convertMetadataToModel(part.Metadata),
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

// convertMetadataToRepo converts metadata values from domain model format to repository model format.
func convertMetadataToRepo(metadata map[string]model.Value) map[string]repoModel.Value {
	result := make(map[string]repoModel.Value)
	for k, v := range metadata {
		var repoValue repoModel.Value

		if v.StringValue != nil {
			repoValue.StringValue = v.StringValue
		}
		if v.Int64Value != nil {
			repoValue.Int64Value = v.Int64Value
		}
		if v.DoubleValue != nil {
			repoValue.DoubleValue = v.DoubleValue
		}
		if v.BoolValue != nil {
			repoValue.BoolValue = v.BoolValue
		}

		result[k] = repoValue
	}
	return result
}

// convertMetadataToModel converts metadata values from repository model format to domain model format.
func convertMetadataToModel(metadata map[string]repoModel.Value) map[string]model.Value {
	result := make(map[string]model.Value)
	for k, v := range metadata {
		var modelValue model.Value

		if v.StringValue != nil {
			modelValue.StringValue = v.StringValue
		}
		if v.Int64Value != nil {
			modelValue.Int64Value = v.Int64Value
		}
		if v.DoubleValue != nil {
			modelValue.DoubleValue = v.DoubleValue
		}
		if v.BoolValue != nil {
			modelValue.BoolValue = v.BoolValue
		}

		result[k] = modelValue
	}
	return result
}
