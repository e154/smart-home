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

package dto

import (
	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Script ...
type Script struct{}

// NewScriptDto ...
func NewScriptDto() Script {
	return Script{}
}

// FromNewScriptRequest ...
func (s Script) FromNewScriptRequest(req *api.NewScriptRequest) (script *m.Script) {
	script = &m.Script{
		Lang:        common.ScriptLang(req.Lang),
		Name:        req.Name,
		Source:      req.Source,
		Description: req.Description,
	}
	return
}

// FromUpdateScriptRequest ...
func (s Script) FromUpdateScriptRequest(req *api.UpdateScriptRequest) (script *m.Script) {
	script = &m.Script{
		Id:          int64(req.Id),
		Lang:        common.ScriptLang(req.Lang),
		Name:        req.Name,
		Source:      req.Source,
		Description: req.Description,
	}
	return
}

// FromExecSrcScriptRequest ...
func (s Script) FromExecSrcScriptRequest(req *api.ExecSrcScriptRequest) (script *m.Script) {
	script = &m.Script{
		Lang:        common.ScriptLang(req.Lang),
		Name:        req.Name,
		Source:      req.Source,
		Description: req.Description,
	}
	return
}

// ToGScript ...
func (s Script) ToGScript(script *m.Script) (result *api.Script) {
	result = &api.Script{
		Id:          int32(script.Id),
		Lang:        string(script.Lang),
		Name:        script.Name,
		Source:      script.Source,
		Description: script.Description,
		CreatedAt:   timestamppb.New(script.CreatedAt),
		UpdatedAt:   timestamppb.New(script.UpdatedAt),
	}
	return
}

// ToSearchResult ...
func (s Script) ToSearchResult(list []*m.Script) *api.SearchScriptListResult {

	items := make([]*api.Script, 0, len(list))

	for _, i := range list {
		items = append(items, s.ToGScript(i))
	}

	return &api.SearchScriptListResult{
		Items: items,
	}
}

// ToListResult ...
func (s Script) ToListResult(list []*m.Script, total, limit, offset uint32) *api.GetScriptListResult {

	items := make([]*api.Script, 0, len(list))

	for _, i := range list {
		items = append(items, s.ToGScript(i))
	}

	return &api.GetScriptListResult{
		Items: items,
		Meta: &api.GetScriptListResult_Meta{
			Limit:        limit,
			ObjectsCount: total,
			Offset:       offset,
		},
	}
}
