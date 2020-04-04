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
	"github.com/e154/smart-home/common"
	"github.com/gin-gonic/gin"
)

type HelloWorld struct {
	appId   string
	session *string
}

func (h HelloWorld) GetAppID() string {
	return h.appId
}

func (h *HelloWorld) OnLaunch(ctx *gin.Context, req *Request, resp *Response) {
	h.session = common.String(req.Session.SessionID)
	resp.OutputSpeech("Hello world from my new Echo test app!").
		Card("Hello World", "This is a test card.").
		EndSession(false)
}

func (h *HelloWorld) OnIntent(ctx *gin.Context, req *Request, resp *Response) {

	if h.session != nil && common.StringValue(h.session) != req.Session.SessionID {
		log.Warn("session expired")
		h.session = nil
		return
	}

	switch req.GetIntentName() {
	default:
		log.Infof("unknown intent name %s", req.GetIntentName())
	}

	resp.OutputSpeech("ok").
		EndSession(true)
}

func (h *HelloWorld) OnSessionEnded(ctx *gin.Context, req *Request, resp *Response) {
	h.session = nil
}

func (h HelloWorld) OnAudioPlayerState(ctx *gin.Context, req *Request, resp *Response) {

}

func NewHelloWorld(appId string) *HelloWorld {
	return &HelloWorld{appId: appId}
}
