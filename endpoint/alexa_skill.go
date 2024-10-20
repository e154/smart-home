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

	"github.com/e154/smart-home/common/apperr"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/alexa"
)

// AlexaSkillEndpoint ...
type AlexaSkillEndpoint struct {
	*CommonEndpoint
}

// NewAlexaSkillEndpoint ...
func NewAlexaSkillEndpoint(common *CommonEndpoint) *AlexaSkillEndpoint {
	return &AlexaSkillEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *AlexaSkillEndpoint) Add(ctx context.Context, params *m.AlexaSkill) (result *m.AlexaSkill, err error) {

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	var id int64
	if id, err = n.adaptors.AlexaSkill.Add(ctx, params); err != nil {
		return
	}

	result, err = n.adaptors.AlexaSkill.GetById(ctx, id)
	if err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/alexa/skill/%d", params.Id), alexa.EventAddedAlexaSkillModel{
		Skill: result,
	})

	log.Infof("added alexa's skill id %d", params.Id)

	return
}

// GetById ...
func (n *AlexaSkillEndpoint) GetById(ctx context.Context, appId int64) (result *m.AlexaSkill, err error) {

	result, err = n.adaptors.AlexaSkill.GetById(ctx, appId)

	return
}

// Update ...
func (n *AlexaSkillEndpoint) Update(ctx context.Context, params *m.AlexaSkill) (skill *m.AlexaSkill, err error) {

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = n.adaptors.AlexaSkill.Update(ctx, params); err != nil {
		return
	}

	skill, err = n.adaptors.AlexaSkill.GetById(ctx, params.Id)
	if err != nil {
		return
	}

	log.Infof("updated alexa's skill id %d", params.Id)

	n.eventBus.Publish(fmt.Sprintf("system/models/alexa/skill/%d", skill.Id), alexa.EventUpdatedAlexaSkillModel{
		Skill: skill,
	})

	return
}

// GetList ...
func (n *AlexaSkillEndpoint) GetList(ctx context.Context, limit, offset int64, order, sortBy string) (result []*m.AlexaSkill, total int64, err error) {

	result, total, err = n.adaptors.AlexaSkill.List(ctx, limit, offset, order, sortBy)

	return
}

// Delete ...
func (n *AlexaSkillEndpoint) Delete(ctx context.Context, skillId int64) (err error) {

	if skillId == 0 {
		err = apperr.ErrInvalidRequest
		return
	}

	var skill *m.AlexaSkill
	skill, err = n.adaptors.AlexaSkill.GetById(ctx, skillId)
	if err != nil {
		return
	}

	if err = n.adaptors.AlexaSkill.Delete(ctx, skill.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/alexa/skill/%d", skillId), alexa.EventDeletedAlexaSkill{
		Skill: skill,
	})

	log.Infof("alexa's skill id %d was deleted", skillId)

	return
}
