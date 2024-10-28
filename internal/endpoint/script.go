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
	"strconv"
	"strings"

	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/models"
	scripts2 "github.com/e154/smart-home/pkg/scripts"

	"github.com/pkg/errors"
)

// ScriptEndpoint ...
type ScriptEndpoint struct {
	*CommonEndpoint
}

// NewScriptEndpoint ...
func NewScriptEndpoint(common *CommonEndpoint) *ScriptEndpoint {
	return &ScriptEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *ScriptEndpoint) Add(ctx context.Context, params *models.Script) (script *models.Script, err error) {

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	var engine scripts2.Engine
	if engine, err = n.scriptService.NewEngine(params); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	if err = engine.Compile(); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	var id int64
	if id, err = n.adaptors.Script.Add(ctx, params); err != nil {
		return
	}

	if script, err = n.adaptors.Script.GetById(ctx, id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/scripts/%d", script.Id), events.EventCreatedScriptModel{
		Common: events.Common{
			Owner: events.OwnerUser,
		},
		ScriptId: script.Id,
		Script:   script,
	})

	log.Infof("added new script %s id:(%d)", params.Name, params.Id)

	return
}

// GetById ...
func (n *ScriptEndpoint) GetById(ctx context.Context, scriptId int64) (result *models.Script, err error) {

	result, err = n.adaptors.Script.GetById(ctx, scriptId)

	return
}

// Copy ...
func (n *ScriptEndpoint) Copy(ctx context.Context, scriptId int64) (script *models.Script, err error) {

	script, err = n.adaptors.Script.GetById(ctx, scriptId)
	if err != nil {
		return
	}

	oldID := script.Id
	oldName := script.Name
	script.Id = 0

	const cpy = "[CPY]"
	if res := strings.Split(script.Name, cpy); len(res) > 1 {
		num, _ := strconv.ParseInt(res[1], 10, 32)
		script.Name = fmt.Sprintf("%s%s%d", res[0], cpy, num+1)
	} else {
		script.Name = fmt.Sprintf("%s%s", script.Name, cpy)
	}

	var id int64
	if id, err = n.adaptors.Script.Add(ctx, script); err != nil {
		return
	}

	if script, err = n.adaptors.Script.GetById(ctx, id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/scripts/%d", script.Id), events.EventCreatedScriptModel{
		Common: events.Common{
			Owner: events.OwnerUser,
		},
		ScriptId: script.Id,
		Script:   script,
	})

	log.Infof("script %s id:(%d) -> %s id:(%d) was copied", oldName, oldID, script.Name, script.Id)

	return
}

// Update ...
func (n *ScriptEndpoint) Update(ctx context.Context, script *models.Script) (result *models.Script, err error) {

	var oldScript *models.Script
	oldScript, err = n.adaptors.Script.GetById(ctx, script.Id)
	if err != nil {
		return
	}

	if ok, errs := n.validation.Valid(script); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	var engine scripts2.Engine
	if engine, err = n.scriptService.NewEngine(script); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	if err = engine.Compile(); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	if err = n.adaptors.Script.Update(ctx, script); err != nil {
		return
	}

	if result, err = n.adaptors.Script.GetById(ctx, script.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/scripts/%d", script.Id), events.EventUpdatedScriptModel{
		Common: events.Common{
			Owner: events.OwnerUser,
		},
		ScriptId:  script.Id,
		Script:    script,
		OldScript: oldScript,
	})

	log.Infof("script %s id:(%d) was updated", script.Name, script.Id)

	return
}

// GetList ...
func (n *ScriptEndpoint) GetList(ctx context.Context, pagination common.PageParams, query *string, ids *[]uint64) (result []*models.Script, total int64, err error) {

	result, total, err = n.adaptors.Script.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, query, ids)

	return
}

// DeleteScriptById ...
func (n *ScriptEndpoint) DeleteScriptById(ctx context.Context, scriptId int64) (err error) {

	if scriptId == 0 {
		err = apperr.ErrBadRequestParams
		return
	}

	var script *models.Script
	script, err = n.adaptors.Script.GetById(ctx, scriptId)
	if err != nil {
		return
	}

	if err = n.adaptors.Script.Delete(ctx, script.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/scripts/%d", script.Id), events.EventRemovedScriptModel{
		Common: events.Common{
			Owner: events.OwnerUser,
		},
		ScriptId: script.Id,
		Script:   script,
	})

	log.Infof("script %s id:(%d) was deleted", script.Name, scriptId)

	return
}

// Execute ...
func (n *ScriptEndpoint) Execute(ctx context.Context, scriptId int64) (result string, err error) {

	var script *models.Script
	script, err = n.adaptors.Script.GetById(ctx, scriptId)
	if err != nil {
		return
	}

	var engine scripts2.Engine
	if engine, err = n.scriptService.NewEngine(script); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	result, err = engine.DoFull()
	if err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
	}

	return
}

// ExecuteSource ...
func (n *ScriptEndpoint) ExecuteSource(ctx context.Context, script *models.Script) (result string, err error) {

	var engine scripts2.Engine
	if engine, err = n.scriptService.NewEngine(script); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	if err = engine.Compile(); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	result, err = engine.DoFull()
	if err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
	}

	return
}

// Search ...
func (n *ScriptEndpoint) Search(ctx context.Context, query string, limit, offset int64) (devices []*models.Script, total int64, err error) {

	devices, total, err = n.adaptors.Script.Search(ctx, query, limit, offset)

	return
}

// Statistic ...
func (n *ScriptEndpoint) Statistic(ctx context.Context) (statistic []*models.Statistic, err error) {
	var stat *models.ScriptsStatistic
	if stat, err = n.adaptors.Script.Statistic(ctx); err != nil {
		return
	}
	statistic = []*models.Statistic{
		{
			Name:        "scripts.stat_total_name",
			Description: "scripts.stat_total_descr",
			Value:       stat.Total,
			Diff:        0,
		},
		{
			Name:        "scripts.stat_used_name",
			Description: "scripts.stat_used_descr",
			Value:       stat.Used,
			Diff:        0,
		},
		{
			Name:        "scripts.stat_unused_name",
			Description: "scripts.stat_unused_descr",
			Value:       stat.Unused,
			Diff:        0,
		},
		{
			Name:        "scripts.stat_js_name",
			Description: "scripts.stat_js_descr",
			Value:       stat.JavaScript,
			Diff:        0,
		},
		{
			Name:        "scripts.stat_cs_name",
			Description: "scripts.stat_cs_descr",
			Value:       stat.CoffeeScript,
			Diff:        0,
		},
		{
			Name:        "scripts.stat_ts_name",
			Description: "scripts.stat_ts_descr",
			Value:       stat.TypeScript,
			Diff:        0,
		},
	}
	return
}
