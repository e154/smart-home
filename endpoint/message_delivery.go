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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/pkg/errors"
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

// GetList ...
func (n *MessageDeliveryEndpoint) GetList(limit, offset int64, order, sortBy string) (result []m.MessageDelivery, total int64, err error) {
	result, total, err = n.adaptors.MessageDelivery.List(limit, offset, order, sortBy)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}

// Delete ...
func (n *MessageDeliveryEndpoint) Delete(id int64) (err error) {
	if err = n.adaptors.MessageDelivery.Delete(id); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}
