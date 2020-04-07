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

package controllers

import (
	"github.com/e154/smart-home/system/stream"
)

type ControllerAction struct {
	*ControllerCommon
}

func NewControllerAction(common *ControllerCommon) *ControllerAction {
	return &ControllerAction{
		ControllerCommon: common,
	}
}

func (c *ControllerAction) Start() {
	c.stream.Subscribe("do.action", c.DoAction)
}

func (c *ControllerAction) Stop() {
	c.stream.UnSubscribe("do.action")
}

// Stream
func (c *ControllerAction) DoAction(client stream.IStreamClient, message stream.Message) {

	v := message.Payload
	var ok bool

	var deviceActionId float64
	var err error

	if deviceActionId, ok = v["action_id"].(float64); !ok {
		log.Warn("bad device_action_id param")
		return
	}

	if _, err = c.core.DoAction(int64(deviceActionId)); err != nil {
		client.Notify("error", err.Error())
		return
	}

	client.Write(message.Success().Pack())
}
