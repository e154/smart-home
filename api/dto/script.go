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

package dto

import (
	stub "github.com/e154/smart-home/api/stub"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// Script ...
type Script struct{}

// NewScriptDto ...
func NewScriptDto() Script {
	return Script{}
}

// FromNewScriptRequest ...
func (s Script) FromNewScriptRequest(req *stub.ApiNewScriptRequest) (script *m.Script) {
	script = &m.Script{
		Lang:        common.ScriptLang(req.Lang),
		Name:        req.Name,
		Source:      req.Source,
		Description: req.Description,
	}
	return
}

// FromUpdateScriptRequest ...
func (s Script) FromUpdateScriptRequest(req *stub.ScriptServiceUpdateScriptByIdJSONBody, id int64) (script *m.Script) {
	script = &m.Script{
		Id:          id,
		Lang:        common.ScriptLang(req.Lang),
		Name:        req.Name,
		Source:      req.Source,
		Description: req.Description,
	}
	return
}

// FromExecSrcScriptRequest ...
func (s Script) FromExecSrcScriptRequest(req *stub.ApiExecSrcScriptRequest) (script *m.Script) {
	script = &m.Script{
		Lang:        common.ScriptLang(req.Lang),
		Name:        req.Name,
		Source:      req.Source,
		Description: req.Description,
	}
	return
}

// GetStubScript ...
func (s Script) GetStubScript(script *m.Script) (result *stub.ApiScript) {
	result = GetStubScript(script)
	return
}

// GetStubScriptShort ...
func (s Script) GetStubScriptShort(script *m.Script) (result *stub.ApiScript) {
	result = GetStubScriptShort(script)
	return
}

// ToSearchResult ...
func (s Script) ToSearchResult(list []*m.Script) *stub.ApiSearchScriptListResult {

	items := make([]stub.ApiScript, 0, len(list))

	for _, script := range list {
		items = append(items, stub.ApiScript{
			Id:   script.Id,
			Lang: string(script.Lang),
			Name: script.Name,
		})
	}

	return &stub.ApiSearchScriptListResult{
		Items: items,
	}
}

// ToListResult ...
func (s Script) ToListResult(list []*m.Script) []*stub.ApiScript {

	items := make([]*stub.ApiScript, 0, len(list))

	for _, script := range list {
		items = append(items, &stub.ApiScript{
			Id:          script.Id,
			Lang:        string(script.Lang),
			Name:        script.Name,
			Description: script.Description,
			CreatedAt:   script.CreatedAt,
			UpdatedAt:   script.UpdatedAt,
		})
	}

	return items
}

// GetStubScript ...
func GetStubScript(script *m.Script) (result *stub.ApiScript) {
	if script == nil {
		return
	}
	result = &stub.ApiScript{
		Id:          script.Id,
		Lang:        string(script.Lang),
		Name:        script.Name,
		Source:      script.Source,
		Description: script.Description,
		ScriptInfo: &stub.ApiScriptInfo{
			AlexaIntents:         int32(script.Info.AlexaIntents),
			EntityActions:        int32(script.Info.EntityActions),
			EntityScripts:        int32(script.Info.EntityScripts),
			AutomationTriggers:   int32(script.Info.AutomationTriggers),
			AutomationConditions: int32(script.Info.AutomationConditions),
			AutomationActions:    int32(script.Info.AutomationActions),
		},
		Versions:  make([]stub.ApiScriptVersion, 0, len(script.Versions)),
		CreatedAt: script.CreatedAt,
		UpdatedAt: script.UpdatedAt,
	}
	for _, version := range script.Versions {
		result.Versions = append(result.Versions, stub.ApiScriptVersion{
			CreatedAt: version.CreatedAt,
			Id:        version.Id,
			Lang:      string(version.Lang),
			Source:    version.Source,
		})
	}
	return
}

// GetStubScriptShort ...
func GetStubScriptShort(script *m.Script) (result *stub.ApiScript) {
	if script == nil {
		return
	}
	result = &stub.ApiScript{
		Id:   script.Id,
		Name: script.Name,
	}
	return
}

func ImportScript(from *stub.ApiScript) (*int64, *m.Script) {
	if from == nil {
		return nil, nil
	}
	return common.Int64(from.Id), &m.Script{
		Id:          from.Id,
		Lang:        common.ScriptLang(from.Lang),
		Name:        from.Name,
		Source:      from.Source,
		Description: from.Description,
	}
}
