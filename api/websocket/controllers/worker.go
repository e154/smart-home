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
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/stream"
)

type ControllerWorker struct {
	*ControllerCommon
}

func NewControllerWorker(common *ControllerCommon) *ControllerWorker {
	return &ControllerWorker{
		ControllerCommon: common,
	}
}

func (c *ControllerWorker) Start() {
	c.stream.Subscribe("do.worker", c.DoWorker)
}

func (c *ControllerWorker) Stop() {
	c.stream.UnSubscribe("do.worker")
}

// Stream
func (c *ControllerWorker) DoWorker(client *stream.Client, message stream.Message) {

	v := message.Payload
	var ok bool

	var workerId float64
	var err error

	if workerId, ok = v["worker_id"].(float64); !ok {
		log.Warning("bad id param")
		return
	}

	var worker *m.Worker
	if worker, err = c.adaptors.Worker.GetById(int64(workerId)); err != nil {
		client.Notify("error", err.Error())
		return
	}

	if err = c.core.DoWorker(worker); err != nil {
		client.Notify("error", err.Error())
		return
	}

	client.Send <- message.Success().Pack()
}
