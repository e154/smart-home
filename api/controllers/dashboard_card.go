// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package controllers

import (
	"context"

	"github.com/e154/smart-home/api/dto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/e154/smart-home/api/stub/api"
)

// ControllerDashboardCard ...
type ControllerDashboardCard struct {
	*ControllerCommon
}

// NewControllerDashboardCard ...
func NewControllerDashboardCard(common *ControllerCommon) ControllerDashboardCard {
	return ControllerDashboardCard{
		ControllerCommon: common,
	}
}

// AddDashboardCard ...
func (c ControllerDashboardCard) AddDashboardCard(ctx context.Context, req *api.NewDashboardCardRequest) (*api.DashboardCard, error) {

	card := c.dto.DashboardCard.AddDashboardCard(req)

	card, errs, err := c.endpoint.DashboardCard.Add(ctx, card)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.DashboardCard.ToDashboardCard(card), nil
}

// UpdateDashboardCard ...
func (c ControllerDashboardCard) UpdateDashboardCard(ctx context.Context, req *api.UpdateDashboardCardRequest) (*api.DashboardCard, error) {

	card := c.dto.DashboardCard.UpdateDashboardCard(req)

	card, errs, err := c.endpoint.DashboardCard.Update(ctx, card)
	if len(errs) != 0 || err != nil {
		return nil, c.error(ctx, errs, err)
	}

	return c.dto.DashboardCard.ToDashboardCard(card), nil
}

// GetDashboardCardById ...
func (c ControllerDashboardCard) GetDashboardCardById(ctx context.Context, req *api.GetDashboardCardRequest) (*api.DashboardCard, error) {

	card, err := c.endpoint.DashboardCard.GetById(ctx, req.Id)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.DashboardCard.ToDashboardCard(card), nil
}

// GetDashboardCardList ...
func (c ControllerDashboardCard) GetDashboardCardList(ctx context.Context, req *api.PaginationRequest) (*api.GetDashboardCardListResult, error) {

	pagination := c.Pagination(req.Page, req.Limit, req.Sort)
	items, total, err := c.endpoint.DashboardCard.GetList(ctx, pagination)
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return c.dto.DashboardCard.ToListResult(items, uint64(total), pagination), nil
}

// DeleteDashboardCard ...
func (c ControllerDashboardCard) DeleteDashboardCard(ctx context.Context, req *api.DeleteDashboardCardRequest) (*emptypb.Empty, error) {

	if err := c.endpoint.DashboardCard.Delete(ctx, req.Id); err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return &emptypb.Empty{}, nil
}

// ImportDashboardCard ...
func (c ControllerDashboardCard) ImportDashboardCard(ctx context.Context, req *api.DashboardCard) (*api.DashboardCard, error) {

	card, err := c.endpoint.DashboardCard.Import(ctx, dto.ImportDashboardCard(req))
	if err != nil {
		return nil, c.error(ctx, nil, err)
	}

	return dto.ToDashboardCard(card), nil
}
