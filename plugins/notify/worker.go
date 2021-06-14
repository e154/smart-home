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

package notify

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"go.uber.org/atomic"
)

// Worker ...
type Worker struct {
	adaptor   *adaptors.Adaptors
	inProcess *atomic.Bool
}

// NewWorker ...
func NewWorker(adaptor *adaptors.Adaptors) *Worker {

	worker := &Worker{
		inProcess: atomic.NewBool(false),
		adaptor:   adaptor,
	}

	return worker
}

func (n *Worker) send(msg m.MessageDelivery, provider Provider) {

	n.inProcess.Store(true)
	defer n.inProcess.Store(false)

	if err := provider.Send(msg.Address, msg.Message); err != nil {
		msg.Status = m.MessageStatusError
		msg.ErrorMessageBody = common.String(err.Error())
	} else {
		msg.Status = m.MessageStatusSucceed
	}
	n.adaptor.MessageDelivery.SetStatus(msg)
}

func (w *Worker) InWork() bool {
	return w.inProcess.Load()
}
