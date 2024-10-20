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

package endpoint

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
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
func (n *Zigbee2mqttEndpoint) AddBridge(ctx context.Context, params *m.Zigbee2mqtt) (bridge *m.Zigbee2mqtt, err error) {

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	params.Id, err = n.adaptors.Zigbee2mqtt.Add(ctx, params)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if bridge, err = n.adaptors.Zigbee2mqtt.GetById(ctx, params.Id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/zigbee2mqtt/%d", bridge.Id), events.EventCreatedZigbee2mqttModel{
		Id:     bridge.Id,
		Bridge: bridge,
	})

	log.Infof("added new z2m %s id:(%d) bridge", bridge.Name, bridge.Id)

	return
}

// GetBridgeById ...
func (n *Zigbee2mqttEndpoint) GetBridgeById(ctx context.Context, id int64) (bridge *m.Zigbee2mqtt, err error) {

	if bridge, err = n.adaptors.Zigbee2mqtt.GetById(ctx, id); err != nil {
		return
	}

	bridge.Info, _ = n.zigbee2mqtt.GetBridgeInfo(id)

	return
}

// UpdateBridge ...
func (n *Zigbee2mqttEndpoint) UpdateBridge(ctx context.Context, params *m.Zigbee2mqtt) (bridge *m.Zigbee2mqtt, err error) {

	if bridge, err = n.adaptors.Zigbee2mqtt.GetById(ctx, params.Id); err != nil {
		return
	}

	bridge.Name = params.Name
	bridge.BaseTopic = params.BaseTopic
	bridge.Login = params.Login
	bridge.Password = params.Password
	bridge.PermitJoin = params.PermitJoin

	if ok, errs := n.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = n.adaptors.Zigbee2mqtt.Update(ctx, bridge); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/zigbee2mqtt/%d", bridge.Id), events.EventUpdatedZigbee2mqttModel{
		Id:     bridge.Id,
		Bridge: bridge,
	})

	log.Infof("updated z2m %s id:(%d) bridge", bridge.Name, bridge.Id)

	return
}

// GetBridgeList ...
func (n *Zigbee2mqttEndpoint) GetBridgeList(ctx context.Context, pagination common.PageParams) (list []*m.Zigbee2mqtt, total int64, err error) {

	if list, total, err = n.adaptors.Zigbee2mqtt.List(ctx, pagination.Limit, pagination.Offset); err != nil {
		return
	}

	for _, br := range list {
		br.Info, _ = n.zigbee2mqtt.GetBridgeInfo(br.Id)
	}

	return
}

// Delete ...
func (n *Zigbee2mqttEndpoint) Delete(ctx context.Context, id int64) (err error) {

	if err = n.adaptors.Zigbee2mqtt.Delete(ctx, id); err != nil {
		return
	}

	n.eventBus.Publish(fmt.Sprintf("system/models/zigbee2mqtt/%d", id), events.EventRemovedZigbee2mqttModel{
		Id: id,
	})

	log.Infof("z2m %d was deleted", id)

	return
}

// ResetBridge ...
func (n *Zigbee2mqttEndpoint) ResetBridge(_ context.Context, id int64) (err error) {

	err = n.zigbee2mqtt.ResetBridge(id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return
		}
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	return
}

// DeviceBan ...
func (n *Zigbee2mqttEndpoint) DeviceBan(_ context.Context, id int64, friendlyName string) (err error) {

	err = n.zigbee2mqtt.BridgeDeviceBan(id, friendlyName)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return
		}
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	return
}

// DeviceWhitelist ...
func (n *Zigbee2mqttEndpoint) DeviceWhitelist(_ context.Context, id int64, friendlyName string) (err error) {

	err = n.zigbee2mqtt.BridgeDeviceWhitelist(id, friendlyName)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return
		}
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	return
}

// Networkmap ...
func (n *Zigbee2mqttEndpoint) Networkmap(_ context.Context, id int64) (networkmap string, err error) {

	networkmap, err = n.zigbee2mqtt.BridgeNetworkmap(id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return
		}
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	return
}

// UpdateNetworkmap ...
func (n *Zigbee2mqttEndpoint) UpdateNetworkmap(_ context.Context, id int64) (err error) {

	err = n.zigbee2mqtt.BridgeUpdateNetworkmap(id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return
		}
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	return
}

// DeviceRename ...
func (n *Zigbee2mqttEndpoint) DeviceRename(_ context.Context, friendlyName, name string) (err error) {

	err = n.zigbee2mqtt.DeviceRename(friendlyName, name)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return
		}
		err = errors.Wrap(apperr.ErrInternal, err.Error())
		return
	}

	return
}

// SearchDevice ...
func (n *Zigbee2mqttEndpoint) SearchDevice(ctx context.Context, search common.SearchParams) (result []*m.Zigbee2mqttDevice, total int64, err error) {

	if result, total, err = n.adaptors.Zigbee2mqttDevice.Search(ctx, search.Query, search.Limit, search.Offset); err != nil {
		err = errors.Wrap(apperr.ErrInternal, err.Error())
	}

	return
}

// DeviceList ...
func (n *Zigbee2mqttEndpoint) DeviceList(ctx context.Context, bridgeId int64, pagination common.PageParams) (result []*m.Zigbee2mqttDevice, total int64, err error) {

	result, total, err = n.adaptors.Zigbee2mqttDevice.ListByBridgeId(ctx, bridgeId, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)

	return
}
