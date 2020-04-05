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

// Javascript Binding
//
// Alexa
//	.OutputSpeech(text)
//	.Card(title, content)
//	.EndSession(bool)
//	.Session()
//

type AlexaBind struct {
	Slots map[string]string
	req   *Request
	resp  *Response
}

func NewAlexaBind(req *Request, resp *Response) (alex *AlexaBind) {
	alex = &AlexaBind{
		Slots: make(map[string]string),
		req:   req,
		resp:  resp,
	}
	for name, slot := range req.Request.Intent.Slots {
		alex.Slots[name] = slot.Value
	}
	return
}

func (r *AlexaBind) OutputSpeech(text string) *AlexaBind {
	r.resp.OutputSpeech(text)
	return r
}

func (r *AlexaBind) Card(title string, content string) *AlexaBind {
	r.resp.Card(title, content)
	return r
}

func (r *AlexaBind) EndSession(flag bool) *AlexaBind {
	r.resp.EndSession(flag)
	return r
}

func (r *AlexaBind) Session() string {
	return r.req.Session.SessionID
}
