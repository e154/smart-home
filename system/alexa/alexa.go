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

type Alexa struct {
}

func NewAlexa() *Alexa {
	return &Alexa{}
}

func echoIntentHandler(echoReq *EchoRequest, echoResp *EchoResponse) {
	echoResp.OutputSpeech("Hello world from my new Echo test app!").Card("Hello World", "This is a test card.")
}

func (a *Alexa) Start() {
	var applications = map[string]interface{}{
		"/echo/helloworld": EchoApplication{ // Route
			AppID:    "amzn1.ask.skill.a8a239d0-3ec9-48b5-9f7d-34c63d9521d5", // Echo App ID from Amazon Dashboard
			OnIntent: echoIntentHandler,
			OnLaunch: echoIntentHandler,
		},
	}

	Run(applications, "3000")
}
