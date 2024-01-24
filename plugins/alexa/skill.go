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

package alexa

import (
	"github.com/pkg/errors"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/gin-gonic/gin"
)

// Skill ...
type Skill struct {
	adaptors      *adaptors.Adaptors
	model         *m.AlexaSkill
	scriptService scripts.ScriptService
	eventBus      bus.Bus
	engine        *scripts.EngineWatcher
	jsBind        *AlexaBind
}

// NewSkill ...
func NewSkill(model *m.AlexaSkill,
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	eventBus bus.Bus) (skill *Skill) {

	skill = &Skill{
		model:         model,
		adaptors:      adaptors,
		scriptService: scriptService,
		eventBus:      eventBus,
		jsBind:        NewAlexaBind(eventBus, model.Id),
	}

	if model.Script != nil {
		skill.engine, _ = scriptService.NewEngineWatcher(model.Script)
		skill.engine.Spawn(func(engine *scripts.Engine) {
			skill.engine.Engine().PushStruct("Alexa", skill.jsBind)
			_, _ = skill.engine.Engine().Do()
		})
	}

	return
}

// GetAppID ...
func (h Skill) GetAppID() string {
	return h.model.SkillId
}

// OnLaunch ...
func (h *Skill) OnLaunch(_ *gin.Context, req *Request, resp *Response) {
	if h.engine == nil {
		return
	}
	h.jsBind.update(req, resp)
	if _, err := h.engine.Engine().AssertFunction("skillOnLaunch"); err != nil {
		log.Error(errors.Wrapf(err, "skill id: %s ", h.model.Id).Error())
	}
}

// OnIntent ...
func (h *Skill) OnIntent(_ *gin.Context, req *Request, resp *Response) {
	var exist bool
	for _, intent := range h.model.Intents {
		if intent.Name != req.GetIntentName() {
			continue
		}
		exist = true

		h.jsBind.update(req, resp)
		if _, err := h.engine.Engine().EvalScript(intent.Script); err != nil {
			log.Error(errors.Wrapf(err, "skill id: %s ", h.model.Id).Error())
			return
		}
		if _, err := h.engine.Engine().AssertFunction("skillOnIntent"); err != nil {
			log.Error(errors.Wrapf(err, "skill id: %s ", h.model.Id).Error())
			return
		}
	}

	if !exist {
		log.Warnf("unknown intent name %s", req.GetIntentName())
	}
}

// OnSessionEnded ...
func (h *Skill) OnSessionEnded(_ *gin.Context, req *Request, resp *Response) {

	if h.engine == nil {
		return
	}

	h.jsBind.update(req, resp)
	if _, err := h.engine.Engine().AssertFunction("skillOnSessionEnd"); err != nil {
		log.Error(errors.Wrapf(err, "skill id: %s ", h.model.Id).Error())
	}
}

// OnAudioPlayerState ...
func (h Skill) OnAudioPlayerState(ctx *gin.Context, req *Request, resp *Response) {

}
