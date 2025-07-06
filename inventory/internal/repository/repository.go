package repository

import (
	"context"

	"github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

type Inventory interface {
	GetPartList(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	AddPart(ctx context.Context, part model.Part) error
	UpdatePart(ctx context.Context, part model.Part) error
	DeletePart(ctx context.Context, uuid string) error
}
