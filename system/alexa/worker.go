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

package alexa

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/scripts"
	"github.com/gin-gonic/gin"
)

// Worker ...
type Worker struct {
	adaptors      *adaptors.Adaptors
	app           *m.AlexaSkill
	scriptService *scripts.ScriptService
	core          *core.Core
}

// NewWorker ...
func NewWorker(app *m.AlexaSkill,
	adaptors *adaptors.Adaptors,
	scriptService *scripts.ScriptService,
	core *core.Core) (worker *Worker) {

	worker = &Worker{
		app:           app,
		adaptors:      adaptors,
		scriptService: scriptService,
		core:          core,
	}

	return
}

// GetAppID ...
func (h Worker) GetAppID() string {
	return h.app.SkillId
}

// OnLaunch ...
func (h *Worker) OnLaunch(ctx *gin.Context, req *Request, resp *Response) {
	if h.app.OnLaunchScript == nil {
		return
	}

	_, err := h.newScript(h.app.OnLaunchScript, req, resp)
	if err != nil {
		log.Error(err.Error())
		return
	}
}

// OnIntent ...
func (h *Worker) OnIntent(ctx *gin.Context, req *Request, resp *Response) {
	var exist bool
	for _, intent := range h.app.Intents {
		if intent.Name != req.GetIntentName() {
			continue
		}
		exist = true

		_, err := h.newScript(intent.Script, req, resp)
		if err != nil {
			log.Error(err.Error())
			return
		}

	}

	if !exist {
		log.Infof("unknown intent name %s", req.GetIntentName())
	}
}

// OnSessionEnded ...
func (h *Worker) OnSessionEnded(ctx *gin.Context, req *Request, resp *Response) {
	if h.app.OnSessionEndScript == nil {
		return
	}

	_, err := h.newScript(h.app.OnSessionEndScript, req, resp)
	if err != nil {
		log.Error(err.Error())
		return
	}
}

// OnAudioPlayerState ...
func (h Worker) OnAudioPlayerState(ctx *gin.Context, req *Request, resp *Response) {

}

func (h *Worker) newScript(s *m.Script, req *Request, resp *Response) (engine *scripts.Engine, err error) {

	if engine, err = h.scriptService.NewEngine(s); err != nil {
		return
	}

	engine.PushStruct("Alexa", NewAlexaBind(req, resp))

	engine.PushFunction("DoAction", h.core.DoAction)

	_, err = engine.Do()

	return
}
