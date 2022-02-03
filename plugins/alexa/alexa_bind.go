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

package alexa

import (
	"github.com/e154/smart-home/system/event_bus"
)

// AlexaBind ...
type AlexaBind struct {
	Slots      map[string]string `json:"slots"`
	req        *Request
	resp       *Response
	eventBus   event_bus.EventBus
	skillId    int64
	intentName string
}

// NewAlexaBind ...
func NewAlexaBind(eventBus event_bus.EventBus, skillId int64) (alex *AlexaBind) {
	alex = &AlexaBind{
		Slots:    make(map[string]string),
		eventBus: eventBus,
		skillId:  skillId,
	}
	return
}

func (r *AlexaBind) update(req *Request, resp *Response) {
	r.req = req
	r.resp = resp
	for name, slot := range req.Request.Intent.Slots {
		r.Slots[name] = slot.Value
	}

	r.intentName = req.Request.Intent.Name
}

// OutputSpeech ...
func (r *AlexaBind) OutputSpeech(text string) *AlexaBind {
	r.resp.OutputSpeech(text)
	return r
}

// Card ...
func (r *AlexaBind) Card(title string, content string) *AlexaBind {
	r.resp.Card(title, content)
	return r
}

// EndSession ...
func (r *AlexaBind) EndSession(flag bool) *AlexaBind {
	r.resp.EndSession(flag)
	return r
}

// Session ...
func (r *AlexaBind) Session() string {
	return r.req.Session.SessionID
}

// SendMessage ...
func (r *AlexaBind) SendMessage(msg interface{}) {
	r.eventBus.Publish(TopicPluginAlexa, EventAlexaAction{
		SkillId:    r.skillId,
		IntentName: r.intentName,
		Payload:    msg,
	})
}
