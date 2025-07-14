package converter

import (
	"github.com/andredubov/rocket-factory/inventory/internal/model"
	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

// PartToRepoModel converts a domain model (repo.Part) to a repository model (repoModel.Part).
func PartToRepoModel(part model.Part) repoModel.Part {
	// Convert metadata from model to repoModel
	metadata := make(map[string]repoModel.Value)
	for k, v := range part.Metadata {
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

		metadata[k] = repoValue
	}

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
		Metadata:  metadata, // Используем сконвертированные метаданные
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

// PartToModel converts a repository model (repoModel.Part) to a domain model (model.Part).
func PartToModel(part repoModel.Part) model.Part {
	// Convert metadata from repoModel to model
	metadata := make(map[string]model.Value)
	for k, v := range part.Metadata {
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

		metadata[k] = modelValue
	}

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
		Metadata:  metadata, // Используем сконвертированные метаданные
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}
