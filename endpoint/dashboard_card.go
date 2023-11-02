// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package endpoint

import (
	"context"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

// DashboardCardEndpoint ...
type DashboardCardEndpoint struct {
	*CommonEndpoint
}

// NewDashboardCardEndpoint ...
func NewDashboardCardEndpoint(common *CommonEndpoint) *DashboardCardEndpoint {
	return &DashboardCardEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (c *DashboardCardEndpoint) Add(ctx context.Context, card *m.DashboardCard) (result *m.DashboardCard, err error) {

	if ok, errs := c.validation.Valid(card); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	var id int64
	if id, err = c.adaptors.DashboardCard.Add(ctx, card); err != nil {
		return
	}

	result, err = c.adaptors.DashboardCard.GetById(ctx, id)

	return
}

// GetById ...
func (c *DashboardCardEndpoint) GetById(ctx context.Context, id int64) (card *m.DashboardCard, err error) {

	if card, err = c.adaptors.DashboardCard.GetById(ctx, id); err != nil {
		return
	}

	err = c.preloadEntities(ctx, card)

	return
}

// Update ...
func (i *DashboardCardEndpoint) Update(ctx context.Context, params *m.DashboardCard) (result *m.DashboardCard, err error) {

	var card *m.DashboardCard
	if card, err = i.adaptors.DashboardCard.GetById(ctx, params.Id); err != nil {
		return
	}

	if err = copier.Copy(&card, &params); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	if ok, errs := i.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = i.adaptors.DashboardCard.Update(ctx, card); err != nil {
		return
	}

	result, err = i.adaptors.DashboardCard.GetById(ctx, params.Id)

	return
}

// GetList ...
func (c *DashboardCardEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []*m.DashboardCard, total int64, err error) {

	list, total, err = c.adaptors.DashboardCard.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	if err != nil {
		return
	}

	for _, card := range list {
		err = c.preloadEntities(ctx, card)
	}

	return
}

// Delete ...
func (c *DashboardCardEndpoint) Delete(ctx context.Context, id int64) (err error) {

	_, err = c.adaptors.DashboardCard.GetById(ctx, id)
	if err != nil {
		return
	}

	err = c.adaptors.DashboardCard.Delete(ctx, id)

	return
}

// Import ...
func (c *DashboardCardEndpoint) Import(ctx context.Context, card *m.DashboardCard) (result *m.DashboardCard, err error) {

	var cardId int64
	if cardId, err = c.adaptors.DashboardCard.Import(ctx, card); err != nil {
		return
	}

	result, err = c.adaptors.DashboardCard.GetById(ctx, cardId)

	return
}

func (c *DashboardCardEndpoint) preloadEntities(ctx context.Context, card *m.DashboardCard) (err error) {

	// get child entities
	entityMap := make(map[common.EntityId]*m.Entity)
	for _, item := range card.Items {
		if item.EntityId != nil {
			entityMap[*item.EntityId] = nil
		}
	}

	entityIds := make([]common.EntityId, 0, len(entityMap))
	for entityId := range entityMap {
		entityIds = append(entityIds, entityId)
	}

	var entites []*m.Entity
	if entites, err = c.adaptors.Entity.GetByIds(ctx, entityIds); err != nil {
		return
	}

	for _, entity := range entites {
		entityMap[entity.Id] = entity
	}

	card.Entities = entityMap

	return
}
