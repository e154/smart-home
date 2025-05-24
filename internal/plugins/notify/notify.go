// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package notify

import (
	"github.com/alitto/pond"
	"github.com/e154/smart-home/internal/plugins/notify/common"
	"github.com/e154/smart-home/pkg/adaptors"
	m "github.com/e154/smart-home/pkg/models"
)

const (
	maxWorkers = 5
	minWorkers = 1
)

type Notify struct {
	adaptors *adaptors.Adaptors
	pool     *pond.WorkerPool
}

// NewNotify ...
func NewNotify(adaptors *adaptors.Adaptors) *Notify {
	return &Notify{
		adaptors: adaptors,
		pool:     pond.New(maxWorkers, 0, pond.MinWorkers(minWorkers)),
	}
}

func (n *Notify) Start() {
	n.pool.RunningWorkers()
}

func (n *Notify) Shutdown() {
	n.pool.StopAndWait()
}

func (n *Notify) Send(msg *m.MessageDelivery, provider Provider) {
	n.pool.Submit(func() {
		NewWorker(n.adaptors).Send(msg, provider)
	})
}

func (n *Notify) SaveAndSend(msg common.Message, provider Provider) {
	n.pool.Submit(func() {
		NewWorker(n.adaptors).SaveAndSend(msg, provider)
	})
}
