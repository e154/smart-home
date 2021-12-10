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
	"github.com/e154/smart-home/system/zigbee2mqtt"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
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

// AddBridge ...
func (n *Zigbee2mqttEndpoint) AddBridge(ctx context.Context, params *m.Zigbee2mqtt) (result *m.Zigbee2mqtt, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = n.validation.Valid(params); !ok {
		return
	}

	if err = n.zigbee2mqtt.AddBridge(params); err != nil {
		return
	}

	result, err = n.zigbee2mqtt.GetBridgeById(params.Id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	return
}

// GetBridgeById ...
func (n *Zigbee2mqttEndpoint) GetBridgeById(ctx context.Context, id int64) (result *zigbee2mqtt.Zigbee2mqttInfo, err error) {

	result, err = n.zigbee2mqtt.GetBridgeInfo(id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}
	return
}

// UpdateBridge ...
func (n *Zigbee2mqttEndpoint) UpdateBridge(ctx context.Context, params *m.Zigbee2mqtt) (bridge *m.Zigbee2mqtt, errs validator.ValidationErrorsTranslations, err error) {

	bridge, err = n.zigbee2mqtt.GetBridgeById(params.Id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	bridge.Password = params.Password
	bridge.BaseTopic = params.BaseTopic
	bridge.Login = params.Login
	bridge.PermitJoin = params.PermitJoin

	var ok bool
	if ok, errs = n.validation.Valid(params); !ok {
		return
	}

	if bridge, err = n.zigbee2mqtt.UpdateBridge(bridge); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}

	return
}

// GetBridgeList ...
func (n *Zigbee2mqttEndpoint) GetBridgeList(ctx context.Context, pagination common.PageParams) (result []*zigbee2mqtt.Zigbee2mqttInfo, total int64, err error) {

	result, total, err = n.zigbee2mqtt.ListBridges(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}

	return
}

// Delete ...
func (n *Zigbee2mqttEndpoint) Delete(ctx context.Context, id int64) (err error) {

	err = n.zigbee2mqtt.DeleteBridge(id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	return
}

// ResetBridge ...
func (n *Zigbee2mqttEndpoint) ResetBridge(ctx context.Context, id int64) (err error) {

	err = n.zigbee2mqtt.ResetBridge(id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	return
}

// DeviceBan ...
func (n *Zigbee2mqttEndpoint) DeviceBan(ctx context.Context, id int64, friendlyName string) (err error) {

	err = n.zigbee2mqtt.BridgeDeviceBan(id, friendlyName)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	return
}

// DeviceWhitelist ...
func (n *Zigbee2mqttEndpoint) DeviceWhitelist(ctx context.Context, id int64, friendlyName string) (err error) {

	err = n.zigbee2mqtt.BridgeDeviceWhitelist(id, friendlyName)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	return
}

// Networkmap ...
func (n *Zigbee2mqttEndpoint) Networkmap(ctx context.Context, id int64) (networkmap string, err error) {

	networkmap, err = n.zigbee2mqtt.BridgeNetworkmap(id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	return
}

// UpdateNetworkmap ...
func (n *Zigbee2mqttEndpoint) UpdateNetworkmap(ctx context.Context, id int64) (err error) {

	err = n.zigbee2mqtt.BridgeUpdateNetworkmap(id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	return
}

// DeviceRename ...
func (n *Zigbee2mqttEndpoint) DeviceRename(ctx context.Context, friendlyName, name string) (err error) {

	err = n.zigbee2mqtt.DeviceRename(friendlyName, name)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	return
}

// SearchDevice ...
func (n *Zigbee2mqttEndpoint) SearchDevice(ctx context.Context, query string, limit, offset int) (result []*m.Zigbee2mqttDevice, total int64, err error) {

	if result, total, err = n.adaptors.Zigbee2mqttDevice.Search(query, limit, offset); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}

	return
}
