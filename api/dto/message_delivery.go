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

package dto

import (
	"fmt"
	stub "github.com/e154/smart-home/api/stub"

	m "github.com/e154/smart-home/models"
)

type MessageDelivery struct{}

func NewMessageDeliveryDto() MessageDelivery {
	return MessageDelivery{}
}

func (m MessageDelivery) ToListResult(list []*m.MessageDelivery) []*stub.ApiMessageDelivery {

	items := make([]*stub.ApiMessageDelivery, 0, len(list))

	for _, i := range list {
		items = append(items, m.ToMessageDelivery(i))
	}

	return items
}

func (m MessageDelivery) ToMessageDelivery(message *m.MessageDelivery) (obj *stub.ApiMessageDelivery) {
	obj = &stub.ApiMessageDelivery{
		Id:                 message.Id,
		Message:            ToMessage(message.Message),
		Address:            message.Address,
		Status:             string(message.Status),
		ErrorMessageStatus: message.ErrorMessageStatus,
		ErrorMessageBody:   message.ErrorMessageBody,
		CreatedAt:          message.CreatedAt,
		UpdatedAt:          message.UpdatedAt,
	}
	return
}

func ToMessage(message *m.Message) (obj stub.ApiMessage) {
	var attributes = make(map[string]string)
	for k, v := range message.Attributes {
		attributes[k] = fmt.Sprintf("%v", v)
	}
	obj = stub.ApiMessage{
		Id:         message.Id,
		Type:       message.Type,
		Attributes: attributes,
		CreatedAt:  message.CreatedAt,
		UpdatedAt:  message.UpdatedAt,
	}
	return
}
