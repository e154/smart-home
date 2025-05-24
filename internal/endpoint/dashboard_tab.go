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

// DashboardTabEndpoint ...
type DashboardTabEndpoint struct {
	*CommonEndpoint
}

// NewDashboardTabEndpoint ...
func NewDashboardTabEndpoint(common *CommonEndpoint) *DashboardTabEndpoint {
	return &DashboardTabEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (t *DashboardTabEndpoint) Add(ctx context.Context, tab *models.DashboardTab) (result *models.DashboardTab, err error) {

	if ok, errs := t.validation.Valid(tab); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	var id int64
	if id, err = t.adaptors.DashboardTab.Add(ctx, tab); err != nil {
		return
	}

	if result, err = t.adaptors.DashboardTab.GetById(ctx, id); err != nil {
		return
	}

	log.Infof("added new dashboard tab %s id:(%d)", result.Name, result.Id)

	return
}

// GetById ...
func (t *DashboardTabEndpoint) GetById(ctx context.Context, id int64) (tab *models.DashboardTab, err error) {

	if tab, err = t.adaptors.DashboardTab.GetById(ctx, id); err != nil {
		return
	}

	err = t.preloadEntities(ctx, tab)

	return
}

// Update ...
func (t *DashboardTabEndpoint) Update(ctx context.Context, params *models.DashboardTab) (result *models.DashboardTab, err error) {

	var tab *models.DashboardTab
	if tab, err = t.adaptors.DashboardTab.GetById(ctx, params.Id); err != nil {
		return
	}

	if err = copier.Copy(&tab, &params); err != nil {
		return
	}

	if ok, errs := t.validation.Valid(tab); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = t.adaptors.DashboardTab.Update(ctx, tab); err != nil {
		return
	}

	if result, err = t.adaptors.DashboardTab.GetById(ctx, params.Id); err != nil {
		return
	}

	log.Infof("updated dashboard tab %s id:(%d)", result.Name, result.Id)

	return
}

// GetList ...
func (t *DashboardTabEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []*models.DashboardTab, total int64, err error) {

	list, total, err = t.adaptors.DashboardTab.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	if err != nil {
		return
	}

	for _, tab := range list {
		err = t.preloadEntities(ctx, tab)
	}

	return
}

// Delete ...
func (t *DashboardTabEndpoint) Delete(ctx context.Context, id int64) (err error) {

	_, err = t.adaptors.DashboardTab.GetById(ctx, id)
	if err != nil {
		return
	}

	if err = t.adaptors.DashboardTab.Delete(ctx, id); err != nil {
		return
	}

	log.Infof("dashboard tab id:(%d) was deleted", id)

	return
}

// Import ...
func (t *DashboardTabEndpoint) Import(ctx context.Context, tab *models.DashboardTab) (result *models.DashboardTab, err error) {

	var id int64
	err = t.adaptors.Transaction.Do(ctx, func(ctx context.Context) error {

		if id, err = t.adaptors.DashboardTab.Add(ctx, tab); err != nil {
			return err
		}

		// cards
		if len(tab.Cards) > 0 {
			for _, card := range tab.Cards {
				card.Id = 0
				card.DashboardTabId = id
				var cardId int64
				if cardId, err = t.adaptors.DashboardCard.Add(ctx, card); err != nil {
					return err
				}

				// items
				if len(card.Items) > 0 {
					for _, item := range card.Items {
						item.Id = 0
						item.DashboardCardId = cardId
						if _, err = t.adaptors.DashboardCardItem.Add(ctx, item); err != nil {
							return err
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

	if result, err = t.adaptors.DashboardTab.GetById(ctx, id); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return
		}
		return
	}

	if err = t.preloadEntities(ctx, tab); err != nil {
		return
	}

	log.Infof("dashboard tab %s id:(%d) was imported", result.Name, result.Id)

	return
}

func (t *DashboardTabEndpoint) preloadEntities(ctx context.Context, tab *models.DashboardTab) (err error) {

	// get child entities
	entityMap := make(map[pkgCommon.EntityId]*models.Entity)
	for _, card := range tab.Cards {
		for _, item := range card.Items {
			if item.EntityId != nil {
				entityMap[*item.EntityId] = nil
			}
		}
	}

	entityIds := make([]pkgCommon.EntityId, 0, len(entityMap))
	for entityId := range entityMap {
		entityIds = append(entityIds, entityId)
	}

	var entities []*models.Entity
	if entities, err = t.adaptors.Entity.GetByIds(ctx, entityIds); err != nil {
		return
	}

	for _, entity := range entities {
		entityMap[entity.Id] = entity
	}

	tab.Entities = entityMap

	return
}
