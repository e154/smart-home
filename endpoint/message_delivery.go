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

package endpoint

import (
	"context"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"strings"
)

// MessageDeliveryEndpoint ...
type MessageDeliveryEndpoint struct {
	*CommonEndpoint
}

// NewMessageDeliveryEndpoint ...
func NewMessageDeliveryEndpoint(common *CommonEndpoint) *MessageDeliveryEndpoint {
	return &MessageDeliveryEndpoint{
		CommonEndpoint: common,
	}
}

// List ...
func (n *MessageDeliveryEndpoint) List(ctx context.Context, pagination common.PageParams, query *string) (result []*m.MessageDelivery, total int64, err error) {
	var messageType []string
	if query != nil {
		messageType = strings.Split(*query, ",")
	}
	result, total, err = n.adaptors.MessageDelivery.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, messageType)
	return
}

// Delete ...
func (n *MessageDeliveryEndpoint) Delete(ctx context.Context, id int64) (err error) {
	err = n.adaptors.MessageDelivery.Delete(ctx, id)
	return
}
