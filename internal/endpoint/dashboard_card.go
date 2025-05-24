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
	"fmt"
	"strings"

	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"

	"github.com/jinzhu/copier"
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
func (c *DashboardCardEndpoint) Add(ctx context.Context, card *models.DashboardCard) (result *models.DashboardCard, err error) {

	if ok, errs := c.validation.Valid(card); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	var id int64
	if id, err = c.adaptors.DashboardCard.Add(ctx, card); err != nil {
		return
	}

	if result, err = c.adaptors.DashboardCard.GetById(ctx, id); err != nil {
		return
	}

	log.Infof("added new dashboard card %s id:(%d)", result.Title, result.Id)

	return
}

// GetById ...
func (c *DashboardCardEndpoint) GetById(ctx context.Context, id int64) (card *models.DashboardCard, err error) {

	if card, err = c.adaptors.DashboardCard.GetById(ctx, id); err != nil {
		return
	}

	err = c.preloadEntities(ctx, card)

	return
}

// Update ...
func (i *DashboardCardEndpoint) Update(ctx context.Context, params *models.DashboardCard) (result *models.DashboardCard, err error) {

	var card *models.DashboardCard
	if card, err = i.adaptors.DashboardCard.GetById(ctx, params.Id); err != nil {
		return
	}

	if err = copier.Copy(&card, &params); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrInternal)
		return
	}

	if ok, errs := i.validation.Valid(params); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	err = i.adaptors.Transaction.Do(ctx, func(ctx context.Context) error {
		if err = i.adaptors.DashboardCard.Update(ctx, card); err != nil {
			return err
		}
		for _, item := range params.Items {
			if err = i.adaptors.DashboardCardItem.Update(ctx, item); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if result, err = i.adaptors.DashboardCard.GetById(ctx, params.Id); err != nil {
		return
	}

	log.Infof("updated dashboard card %s id:(%d)", result.Title, result.Id)

	return
}

// GetList ...
func (c *DashboardCardEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []*models.DashboardCard, total int64, err error) {

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

	if err = c.adaptors.DashboardCard.Delete(ctx, id); err != nil {
		return
	}

	log.Infof("dashboard card id:(%d) was deleted", id)

	return
}

// Import ...
func (c *DashboardCardEndpoint) Import(ctx context.Context, card *models.DashboardCard) (result *models.DashboardCard, err error) {

	var cardId int64
	err = c.adaptors.Transaction.Do(ctx, func(ctx context.Context) error {

		card.Id = 0

		if !strings.Contains(card.Title, "[IMPORTED]") {
			card.Title = card.Title + " [IMPORTED]"
		}

		if cardId, err = c.adaptors.DashboardCard.Add(ctx, card); err != nil {
			return err
		}

		// items
		if len(card.Items) > 0 {
			for _, item := range card.Items {
				item.Id = 0
				item.DashboardCardId = cardId
				if _, err = c.adaptors.DashboardCardItem.Add(ctx, item); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if result, err = c.adaptors.DashboardCard.GetById(ctx, cardId); err != nil {
		return
	}

	log.Infof("dashboard card %s id:(%d) was imported", result.Title, result.Id)

	return
}

func (c *DashboardCardEndpoint) preloadEntities(ctx context.Context, card *models.DashboardCard) (err error) {

	// get child entities
	entityMap := make(map[pkgCommon.EntityId]*models.Entity)
	for _, item := range card.Items {
		if item.EntityId != nil {
			entityMap[*item.EntityId] = nil
		}
	}

	entityIds := make([]pkgCommon.EntityId, 0, len(entityMap))
	for entityId := range entityMap {
		entityIds = append(entityIds, entityId)
	}

	var entites []*models.Entity
	if entites, err = c.adaptors.Entity.GetByIds(ctx, entityIds); err != nil {
		return
	}

	for _, entity := range entites {
		entityMap[entity.Id] = entity
	}

	card.Entities = entityMap

	return
}
