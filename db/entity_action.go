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

package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// EntityActions ...
type EntityActions struct {
	Db *gorm.DB
}

// EntityAction ...
type EntityAction struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Icon        *string
	Entity      *Entity
	EntityId    common.EntityId
	Image       *Image
	ImageId     *int64
	Script      *Script
	ScriptId    *int64
	Type        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *EntityAction) TableName() string {
	return "entity_actions"
}

// Add ...
func (n EntityActions) Add(ctx context.Context, v *EntityAction) (id int64, err error) {
	if err = n.Db.WithContext(ctx).Create(&v).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "name_at_entity_actions_unq") {
					err = errors.Wrap(apperr.ErrEntityActionAdd, fmt.Sprintf("action name \"%s\" not unique", v.Name))
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrEntityActionAdd, err.Error())
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n EntityActions) GetById(ctx context.Context, id int64) (v *EntityAction, err error) {
	v = &EntityAction{Id: id}
	if err = n.Db.WithContext(ctx).First(&v).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrEntityActionNotFound, fmt.Sprintf("id \"%d\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrEntityActionGet, err.Error())
	}
	return
}

// Update ...
func (n EntityActions) Update(ctx context.Context, m *EntityAction) (err error) {
	err = n.Db.WithContext(ctx).Model(&EntityAction{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"icon":        m.Icon,
		"entity_id":   m.EntityId,
		"image_id":    m.ImageId,
		"script_id":   m.ScriptId,
		"type":        m.Type,
	}).Error

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "name_at_entity_actions_unq") {
					err = errors.Wrap(apperr.ErrEntityActionUpdate, fmt.Sprintf("action name \"%s\" not unique", m.Name))
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrEntityActionUpdate, err.Error())
	}
	return
}

// DeleteByEntityId ...
func (n EntityActions) DeleteByEntityId(ctx context.Context, deviceId common.EntityId) (err error) {
	if err = n.Db.WithContext(ctx).Delete(&EntityAction{}, "entity_id = ?", deviceId).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityActionDelete, err.Error())
	}
	return
}

// List ...
func (n *EntityActions) List(ctx context.Context, limit, offset int, orderBy, sort string) (list []*EntityAction, total int64, err error) {

	if err = n.Db.WithContext(ctx).Model(EntityAction{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrEntityActionList, err.Error())
		return
	}

	list = make([]*EntityAction, 0)
	err = n.Db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	if err != nil {
		err = errors.Wrap(apperr.ErrEntityActionList, err.Error())
	}
	return
}

// AddMultiple ...
func (n *EntityActions) AddMultiple(ctx context.Context, actions []*EntityAction) (err error) {
	if err = n.Db.WithContext(ctx).Create(&actions).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "name_at_entity_states_unq") {
					err = errors.Wrap(apperr.ErrEntityActionAdd, "multiple insert")
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrEntityActionAdd, err.Error())
	}
	return
}
