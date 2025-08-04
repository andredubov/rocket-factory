package inventory

import (
	"context"

	"github.com/andredubov/rocket-factory/order/internal/client/converter"
	"github.com/andredubov/rocket-factory/order/internal/model"
)

// ListParts retrieves a list of parts from the inventory service based on the provided filter.
func (c *inventoryClient) ListParts(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	request := converter.PartFilterToListPartsRequest(filter)

	response, err := c.generatedClient.ListParts(ctx, request)
	if err != nil {
		return nil, err
	}

	return converter.PartsFromListPartsResponse(response), nil
}
