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
	"errors"

	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"
	"github.com/jinzhu/copier"
)

// DashboardEndpoint ...
type DashboardEndpoint struct {
	*CommonEndpoint
}

// NewDashboardEndpoint ...
func NewDashboardEndpoint(common *CommonEndpoint) *DashboardEndpoint {
	return &DashboardEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (d *DashboardEndpoint) Add(ctx context.Context, board *models.Dashboard) (result *models.Dashboard, err error) {

	if ok, errs := d.validation.Valid(board); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	var id int64
	if id, err = d.adaptors.Dashboard.Add(ctx, board); err != nil {
		return
	}

	if result, err = d.adaptors.Dashboard.GetById(ctx, id); err != nil {
		return nil, err
	}

	log.Infof("added new dashboard %s id:(%d)", result.Name, result.Id)

	return
}

// GetById ...
func (d *DashboardEndpoint) GetById(ctx context.Context, id int64) (board *models.Dashboard, err error) {

	if board, err = d.adaptors.Dashboard.GetById(ctx, id); err != nil {
		return
	}

	err = d.preloadEntities(ctx, board)

	return
}

// Search ...
func (d *DashboardEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*models.Dashboard, total int64, err error) {

	if limit == 0 {
		limit = common.DefaultPageSize
	}

	result, total, err = d.adaptors.Dashboard.Search(ctx, query, limit, offset)

	return
}

// Update ...
func (i *DashboardEndpoint) Update(ctx context.Context, params *models.Dashboard) (result *models.Dashboard, err error) {

	var board *models.Dashboard
	if board, err = i.adaptors.Dashboard.GetById(ctx, params.Id); err != nil {
		return
	}

	if err = copier.Copy(&board, &params); err != nil {
		return
	}

	if ok, errs := i.validation.Valid(params); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = i.adaptors.Dashboard.Update(ctx, board); err != nil {
		return
	}

	if result, err = i.adaptors.Dashboard.GetById(ctx, params.Id); err != nil {
		if !errors.Is(err, apperr.ErrNotFound) {
		}
	}

	log.Infof("updated dashboard %s id:(%d)", result.Name, result.Id)

	return
}

// GetList ...
func (d *DashboardEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []*models.Dashboard, total int64, err error) {

	list, total, err = d.adaptors.Dashboard.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	if err != nil {
		return
	}

	return
}

// Delete ...
func (d *DashboardEndpoint) Delete(ctx context.Context, id int64) (err error) {

	_, err = d.adaptors.Dashboard.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return
		}
		return
	}

	if err = d.adaptors.Dashboard.Delete(ctx, id); err != nil {
		return err
	}

	log.Infof("dashboard id:(%d) was deleted", id)

	return
}

func (d *DashboardEndpoint) preloadEntities(ctx context.Context, board *models.Dashboard) (err error) {

	// get child entities
	entityMap := make(map[pkgCommon.EntityId]*models.Entity)
	for _, tab := range board.Tabs {
		for _, card := range tab.Cards {
			if card.EntityId != nil {
				entityMap[*card.EntityId] = nil
			}
			for _, item := range card.Items {
				if item.EntityId != nil {
					entityMap[*item.EntityId] = nil
				}
			}
		}
	}

	entityIds := make([]pkgCommon.EntityId, 0, len(entityMap))
	for entityId := range entityMap {
		entityIds = append(entityIds, entityId)
	}

	var entities []*models.Entity
	if entities, err = d.adaptors.Entity.GetByIds(ctx, entityIds); err != nil {
		return
	}

	for _, entity := range entities {
		entityMap[entity.Id] = entity
	}

	board.Entities = entityMap

	return
}

// Import ...
func (d *DashboardEndpoint) Import(ctx context.Context, board *models.Dashboard) (result *models.Dashboard, err error) {

	board.Id = 0
	board.Name = board.Name + " [IMPORTED]"
	var id int64

	err = d.adaptors.Transaction.Do(ctx, func(ctx context.Context) error {
		if id, err = d.adaptors.Dashboard.Add(ctx, board); err != nil {
			return err
		}

		// tabs
		if len(board.Tabs) > 0 {
			for _, tab := range board.Tabs {
				tab.Id = 0
				tab.DashboardId = id
				var tabId int64
				if tabId, err = d.adaptors.DashboardTab.Add(ctx, tab); err != nil {
					return err
				}

				// cards
				if len(tab.Cards) > 0 {
					for _, card := range tab.Cards {
						card.Id = 0
						card.DashboardTabId = tabId
						var cardId int64
						if cardId, err = d.adaptors.DashboardCard.Add(ctx, card); err != nil {
							return err
						}

						// items
						if len(card.Items) > 0 {
							for _, item := range card.Items {
								item.Id = 0
								item.DashboardCardId = cardId
								if _, err = d.adaptors.DashboardCardItem.Add(ctx, item); err != nil {
									return err
								}
							}
						}
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if result, err = d.adaptors.Dashboard.GetById(ctx, id); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return
		}
		return
	}

	if err = d.preloadEntities(ctx, board); err != nil {
		return nil, err
	}

	log.Infof("dashboard %s id:(%d) was imported", result.Name, id)

	return
}
