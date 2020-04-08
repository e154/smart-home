// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"errors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
)

type AlexaSkillEndpoint struct {
	*CommonEndpoint
}

func NewAlexaApplicationEndpoint(common *CommonEndpoint) *AlexaSkillEndpoint {
	return &AlexaSkillEndpoint{
		CommonEndpoint: common,
	}
}

func (n *AlexaSkillEndpoint) Add(params *m.AlexaSkill) (result *m.AlexaSkill, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = n.adaptors.AlexaSkill.Add(params); err != nil {
		return
	}

	if result, err = n.adaptors.AlexaSkill.GetById(id); err != nil {
		return
	}

	n.alexa.Add(result)

	return
}

func (n *AlexaSkillEndpoint) GetById(appId int64) (result *m.AlexaSkill, err error) {

	result, err = n.adaptors.AlexaSkill.GetById(appId)

	return
}

func (n *AlexaSkillEndpoint) Update(params *m.AlexaSkill) (skill *m.AlexaSkill, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.AlexaSkill.Update(params); err != nil {
		return
	}

	if skill, err = n.adaptors.AlexaSkill.GetById(params.Id); err != nil {
		return
	}

	n.alexa.Update(skill)

	return
}

func (n *AlexaSkillEndpoint) GetList(limit, offset int64, order, sortBy string) (result []*m.AlexaSkill, total int64, err error) {

	result, total, err = n.adaptors.AlexaSkill.List(limit, offset, order, sortBy)

	return
}

func (n *AlexaSkillEndpoint) Delete(skillId int64) (err error) {

	if skillId == 0 {
		err = errors.New("skill id is null")
		return
	}

	var skill *m.AlexaSkill
	if skill, err = n.adaptors.AlexaSkill.GetById(skillId); err != nil {
		return
	}

	if err = n.adaptors.AlexaSkill.Delete(skill.Id); err != nil {
		return
	}

	n.alexa.Delete(skill)

	return
}
