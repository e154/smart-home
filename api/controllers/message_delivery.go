package controllers

import (
	"context"

	"github.com/e154/smart-home/api/stub/api"
)

// ControllerMessageDelivery ...
type ControllerMessageDelivery struct {
	*ControllerCommon
}

func NewControllerMessageDelivery(common *ControllerCommon) ControllerMessageDelivery {
	return ControllerMessageDelivery{
		ControllerCommon: common,
	}
}

// GetMessageDeliveryList ...
func (c ControllerMessageDelivery) GetMessageDeliveryList(ctx context.Context, req *api.MessageDeliveryPaginationRequest) (*api.GetMessageDeliveryListResult, error) {
	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.MessageDelivery.List(ctx, pagination, req.MessageType, req.StartDate, req.EndDate)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.MessageDelivery.ToListResult(items, uint64(total), pagination), nil
}
