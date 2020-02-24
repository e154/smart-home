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

package endpoint

import (
	"errors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
)

type Zigbee2mqttEndpoint struct {
	*CommonEndpoint
}

func NewZigbee2mqttEndpoint(common *CommonEndpoint) *Zigbee2mqttEndpoint {
	return &Zigbee2mqttEndpoint{
		CommonEndpoint: common,
	}
}

func (n *Zigbee2mqttEndpoint) Add(params *m.Zigbee2mqtt) (result *m.Zigbee2mqtt, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.zigbee2mqtt.AddBridge(params); err != nil {
		return
	}

	result, err = n.zigbee2mqtt.GetBridgeById(params.Id)

	return
}

func (n *Zigbee2mqttEndpoint) GetById(id int64) (result *m.Zigbee2mqtt, err error) {

	result, err = n.zigbee2mqtt.GetBridgeById(id)

	return
}

func (n *Zigbee2mqttEndpoint) Update(params *m.Zigbee2mqtt) (bridge *m.Zigbee2mqtt, errs []*validation.Error, err error) {

	bridge, err = n.zigbee2mqtt.GetBridgeById(params.Id)
	if err != nil {
		return
	}

	common.Copy(&bridge, &params, common.JsonEngine)

	// validation
	_, errs = bridge.Valid()
	if len(errs) > 0 {
		return
	}

	bridge, err = n.zigbee2mqtt.UpdateBridge(bridge)

	return
}

func (n *Zigbee2mqttEndpoint) GetList(limit, offset int64, order, sortBy string) (result []m.Zigbee2mqtt, total int64, err error) {

	result, total, err = n.zigbee2mqtt.ListBridges(limit, offset, order, sortBy)

	return
}

func (n *Zigbee2mqttEndpoint) Delete(id int64) (err error) {

	if id == 0 {
		err = errors.New("node id is null")
		return
	}

	err = n.zigbee2mqtt.DeleteBridge(id)

	return
}
