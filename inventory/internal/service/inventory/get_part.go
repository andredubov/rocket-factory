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

func (i *inventoryService) GetPart(ctx context.Context, partUUID string) (*model.Part, error) {
	ctx, span := tracing.StartSpan(ctx, "inventory.get_part",
		trace.WithAttributes(
			attribute.String("part.uuid", partUUID),
		),
	)
	defer span.End()

	part, err := i.inventoryRepository.GetPart(ctx, partUUID)
	if err != nil {
		logger.Error(ctx, "❌ [InventoryService] failed to get part by its uuid from database", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.String("part.uuid", part.Uuid))
	span.SetStatus(codes.Ok, "get part succeeded")

	return part, nil
}
