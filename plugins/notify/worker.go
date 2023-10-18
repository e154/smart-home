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

package notify

import (
	"context"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// Worker ...
type Worker struct {
	adaptor *adaptors.Adaptors
}

// NewWorker ...
func NewWorker(adaptor *adaptors.Adaptors) *Worker {
	return &Worker{
		adaptor: adaptor,
	}
}

func (n *Worker) SaveAndSend(msg Message, provider Provider) {
	addresses, message := provider.Save(msg)

	var err error

	for _, address := range addresses {
		messageDelivery := &m.MessageDelivery{
			Message:   message,
			MessageId: message.Id,
			EntityId:  msg.EntityId,
			Status:    m.MessageStatusInProgress,
			Address:   address,
		}
		if messageDelivery.Id, err = n.adaptor.MessageDelivery.Add(context.Background(), messageDelivery); err != nil {
			log.Error(err.Error())
		}
		n.Send(messageDelivery, provider)
	}
}

func (n *Worker) Send(msg *m.MessageDelivery, provider Provider) {

	if err := provider.Send(msg.Address, msg.Message); err != nil {
		msg.Status = m.MessageStatusError
		msg.ErrorMessageBody = common.String(err.Error())
	} else {
		msg.Status = m.MessageStatusSucceed
	}
	_ = n.adaptor.MessageDelivery.SetStatus(context.Background(), msg)
}
