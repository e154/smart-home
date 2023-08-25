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

package endpoint

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
	"github.com/go-playground/validator/v10"
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
func (n *ScriptEndpoint) Add(ctx context.Context, params *m.Script) (result *m.Script, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = n.validation.Valid(params); !ok {
		return
	}

	var engine *scripts.Engine
	if engine, err = n.scriptService.NewEngine(params); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	if err = engine.Compile(); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	var id int64
	if id, err = n.adaptors.Script.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.Script.GetById(id)

	return
}

// GetById ...
func (n *ScriptEndpoint) GetById(ctx context.Context, scriptId int64) (result *m.Script, err error) {

	result, err = n.adaptors.Script.GetById(scriptId)

	return
}

// Copy ...
func (n *ScriptEndpoint) Copy(ctx context.Context, scriptId int64) (script *m.Script, err error) {

	script, err = n.adaptors.Script.GetById(scriptId)
	if err != nil {
		return
	}

	script.Id = 0

	const cpy = "[CPY]"
	if res := strings.Split(script.Name, cpy); len(res) > 1 {
		num, _ := strconv.ParseInt(res[1], 10, 32)
		script.Name = fmt.Sprintf("%s%s%d", res[0], cpy, num+1)
	} else {
		script.Name = fmt.Sprintf("%s%s", script.Name, cpy)
	}

	var id int64
	if id, err = n.adaptors.Script.Add(script); err != nil {
		return
	}

	script, err = n.adaptors.Script.GetById(id)

	return
}

// Update ...
func (n *ScriptEndpoint) Update(ctx context.Context, params *m.Script) (result *m.Script, errs validator.ValidationErrorsTranslations, err error) {

	var script *m.Script
	script, err = n.adaptors.Script.GetById(params.Id)
	if err != nil {
		return
	}

	if err = common.Copy(&script, &params); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	var ok bool
	if ok, errs = n.validation.Valid(params); !ok {
		return
	}

	var engine *scripts.Engine
	if engine, err = n.scriptService.NewEngine(script); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	if err = engine.Compile(); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	if err = n.adaptors.Script.Update(script); err != nil {
		return
	}

	result, err = n.adaptors.Script.GetById(script.Id)

	return
}

// GetList ...
func (n *ScriptEndpoint) GetList(ctx context.Context, pagination common.PageParams) (result []*m.Script, total int64, err error) {

	result, total, err = n.adaptors.Script.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)

	return
}

// DeleteScriptById ...
func (n *ScriptEndpoint) DeleteScriptById(ctx context.Context, scriptId int64) (err error) {

	if scriptId == 0 {
		err = apperr.ErrBadRequestParams
		return
	}

	var script *m.Script
	script, err = n.adaptors.Script.GetById(scriptId)
	if err != nil {
		return
	}

	err = n.adaptors.Script.Delete(script.Id)

	return
}

// Execute ...
func (n *ScriptEndpoint) Execute(ctx context.Context, scriptId int64) (result string, err error) {

	var script *m.Script
	script, err = n.adaptors.Script.GetById(scriptId)
	if err != nil {
		return
	}

	var engine *scripts.Engine
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
func (n *ScriptEndpoint) ExecuteSource(ctx context.Context, script *m.Script) (result string, err error) {

	var engine *scripts.Engine
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
func (n *ScriptEndpoint) Search(ctx context.Context, query string, limit, offset int64) (devices []*m.Script, total int64, err error) {

	devices, total, err = n.adaptors.Script.Search(query, limit, offset)

	return
}

// Statistic ...
func (n *ScriptEndpoint) Statistic(ctx context.Context) (statistic []*m.Statistic, err error) {
	var stat *m.ScriptsStatistic
	if stat, err = n.adaptors.Script.Statistic(); err != nil {
		return
	}
	statistic = []*m.Statistic{
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
