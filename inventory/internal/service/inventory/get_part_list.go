package inventory

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/andredubov/rocket-factory/inventory/internal/model"
	"github.com/andredubov/rocket-factory/platform/pkg/logger"
	"github.com/andredubov/rocket-factory/platform/pkg/tracing"
)

func (i *inventoryService) GetPartList(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	ctx, span := tracing.StartSpan(ctx, "inventory.get_part_list",
		trace.WithAttributes(
			attribute.Int("filter.UUIDs.count", len(filter.UUIDs)),
		),
	)
	defer span.End()

	parts, err := i.inventoryRepository.GetPartList(ctx, filter)
	if err != nil {
		logger.Error(ctx, "❌ [InventoryService] failed to get part list from database", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "get part list succeeded")

	return parts, nil
}
