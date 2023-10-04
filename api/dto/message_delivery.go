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

	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageDelivery struct{}

func NewMessageDeliveryDto() MessageDelivery {
	return MessageDelivery{}
}

func (m MessageDelivery) ToListResult(list []*m.MessageDelivery, total uint64, pagination common.PageParams) *api.GetMessageDeliveryListResult {

	items := make([]*api.MessageDelivery, 0, len(list))

	for _, i := range list {
		items = append(items, m.ToMessageDelivery(i))
	}

	return &api.GetMessageDeliveryListResult{
		Items: items,
		Meta: &api.Meta{
			Limit: uint64(pagination.Limit),
			Page:  pagination.PageReq,
			Total: total,
			Sort:  pagination.SortReq,
		},
	}
}

func (m MessageDelivery) ToMessageDelivery(message *m.MessageDelivery) (obj *api.MessageDelivery) {
	obj = &api.MessageDelivery{
		Id:                 message.Id,
		Message:            ToMessage(message.Message),
		Address:            message.Address,
		Status:             string(message.Status),
		ErrorMessageStatus: message.ErrorMessageStatus,
		ErrorMessageBody:   message.ErrorMessageBody,
		CreatedAt:          timestamppb.New(message.CreatedAt),
		UpdatedAt:          timestamppb.New(message.UpdatedAt),
	}
	return
}

func ToMessage(message *m.Message) (obj *api.Message) {
	var attributes = make(map[string]string)
	for k, v := range message.Attributes {
		attributes[k] = fmt.Sprintf("%v", v)
	}
	obj = &api.Message{
		Id:         message.Id,
		Type:       message.Type,
		Attributes: attributes,
		CreatedAt:  timestamppb.New(message.CreatedAt),
		UpdatedAt:  timestamppb.New(message.UpdatedAt),
	}
	return
}
