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
	"errors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/system/zigbee2mqtt"
)

// Zigbee2mqttEndpoint ...
type Zigbee2mqttEndpoint struct {
	*CommonEndpoint
}

// NewZigbee2mqttEndpoint ...
func NewZigbee2mqttEndpoint(common *CommonEndpoint) *Zigbee2mqttEndpoint {
	return &Zigbee2mqttEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
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

// GetById ...
func (n *Zigbee2mqttEndpoint) GetById(id int64) (result *zigbee2mqtt.Zigbee2mqttInfo, err error) {

	result, err = n.zigbee2mqtt.GetBridgeInfo(id)

	return
}

// Update ...
func (n *Zigbee2mqttEndpoint) Update(params *m.Zigbee2mqtt) (bridge *m.Zigbee2mqtt, errs []*validation.Error, err error) {

	bridge, err = n.zigbee2mqtt.GetBridgeById(params.Id)
	if err != nil {
		return
	}

	bridge.Password = params.Password
	bridge.BaseTopic = params.BaseTopic
	bridge.Login = params.Login
	bridge.PermitJoin = params.PermitJoin

	// validation
	_, errs = bridge.Valid()
	if len(errs) > 0 {
		return
	}

	bridge, err = n.zigbee2mqtt.UpdateBridge(bridge)

	return
}

// GetList ...
func (n *Zigbee2mqttEndpoint) GetList(limit, offset int64, order, sortBy string) (result []*zigbee2mqtt.Zigbee2mqttInfo, total int64, err error) {

	result, total, err = n.zigbee2mqtt.ListBridges(limit, offset, order, sortBy)

	return
}

// Delete ...
func (n *Zigbee2mqttEndpoint) Delete(id int64) (err error) {

	if id == 0 {
		err = errors.New("node id is null")
		return
	}

	err = n.zigbee2mqtt.DeleteBridge(id)

	return
}

// ResetBridge ...
func (n *Zigbee2mqttEndpoint) ResetBridge(id int64) (err error) {

	err = n.zigbee2mqtt.ResetBridge(id)

	return
}

// DeviceBan ...
func (n *Zigbee2mqttEndpoint) DeviceBan(id int64, friendlyName string) (err error) {

	err = n.zigbee2mqtt.BridgeDeviceBan(id, friendlyName)

	return
}

// DeviceWhitelist ...
func (n *Zigbee2mqttEndpoint) DeviceWhitelist(id int64, friendlyName string) (err error) {

	err = n.zigbee2mqtt.BridgeDeviceWhitelist(id, friendlyName)

	return
}

// Networkmap ...
func (n *Zigbee2mqttEndpoint) Networkmap(id int64) (networkmap string, err error) {

	networkmap, err = n.zigbee2mqtt.BridgeNetworkmap(id)

	return
}

// UpdateNetworkmap ...
func (n *Zigbee2mqttEndpoint) UpdateNetworkmap(id int64) (err error) {

	err = n.zigbee2mqtt.BridgeUpdateNetworkmap(id)

	return
}

// DeviceRename ...
func (n *Zigbee2mqttEndpoint) DeviceRename(friendlyName, name string) (err error) {

	err = n.zigbee2mqtt.DeviceRename(friendlyName, name)

	return
}

// SearchDevice ...
func (n *Zigbee2mqttEndpoint) SearchDevice(query string, limit, offset int) (result []*m.Zigbee2mqttDevice, total int64, err error) {

	result, total, err = n.adaptors.Zigbee2mqttDevice.Search(query, limit, offset)

	return
}
