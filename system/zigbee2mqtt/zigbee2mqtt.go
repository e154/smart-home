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

package zigbee2mqtt

import (
	"context"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/op/go-logging"
	"sync"
)

var (
	log = logging.MustGetLogger("zigbee2mqtt")
)

type Zigbee2mqtt struct {
	metric      *metrics.MetricManager
	graceful    *graceful_service.GracefulService
	mqtt        *mqtt.Mqtt
	adaptors    *adaptors.Adaptors
	isStarted   bool
	bridgesLock *sync.Mutex
	bridges     map[int64]*Bridge
}

func NewZigbee2mqtt(graceful *graceful_service.GracefulService,
	mqtt *mqtt.Mqtt,
	adaptors *adaptors.Adaptors,
	metric *metrics.MetricManager) *Zigbee2mqtt {
	return &Zigbee2mqtt{
		graceful:    graceful,
		mqtt:        mqtt,
		adaptors:    adaptors,
		bridgesLock: &sync.Mutex{},
		bridges:     make(map[int64]*Bridge),
		metric:      metric,
	}
}

func (z *Zigbee2mqtt) Start() {
	if z.isStarted {
		return
	}
	z.isStarted = true

	models, _, err := z.adaptors.Zigbee2mqtt.List(99, 0)
	if err != nil {
		log.Error(err.Error())
	}

	if len(models) == 0 {
		model := &m.Zigbee2mqtt{
			Name:       "zigbee2mqtt",
			BaseTopic:  "zigbee2mqtt",
			PermitJoin: true,
		}
		model.Id, err = z.adaptors.Zigbee2mqtt.Add(model)
		if err != nil {
			log.Error(err.Error())
			return
		}
		models = append(models, model)
	}

	for _, model := range models {
		bridge := NewBridge(z.mqtt, z.adaptors, model, z.metric)
		bridge.Start()

		z.bridgesLock.Lock()
		z.bridges[model.Id] = bridge
		z.bridgesLock.Unlock()
	}
}

func (z *Zigbee2mqtt) Shutdown() {
	if !z.isStarted {
		return
	}
	z.isStarted = false
	for _, bridge := range z.bridges {
		bridge.Stop(context.Background())
	}
}

func (z *Zigbee2mqtt) AddBridge(model *m.Zigbee2mqtt) (err error) {

	model.Id, err = z.adaptors.Zigbee2mqtt.Add(model)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if model, err = z.adaptors.Zigbee2mqtt.GetById(model.Id); err != nil {
		return
	}

	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	bridge := NewBridge(z.mqtt, z.adaptors, model, z.metric)
	bridge.Start()
	z.bridges[model.Id] = bridge
	return
}

func (z *Zigbee2mqtt) GetBridgeById(id int64) (*m.Zigbee2mqtt, error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	if br, ok := z.bridges[id]; ok {
		model := br.GetModel()
		return &model, nil
	}
	return nil, adaptors.ErrRecordNotFound
}

func (z *Zigbee2mqtt) GetBridgeInfo(id int64) (*Zigbee2mqttInfo, error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	if br, ok := z.bridges[id]; ok {
		return br.Info(), nil
	}
	return nil, adaptors.ErrRecordNotFound
}

func (z *Zigbee2mqtt) ListBridges(limit, offset int64, order, sortBy string) (models []*Zigbee2mqttInfo, total int64, err error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	total = int64(len(z.bridges))

	for _, br := range z.bridges {
		models = append(models, br.Info())
	}

	return
}

func (z *Zigbee2mqtt) UpdateBridge(model *m.Zigbee2mqtt) (result *m.Zigbee2mqtt, err error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	var bridge *Bridge
	if bridge, err = z.unsafeGetBridge(model.Id); err == nil {

	} else {
		return
	}

	if err = z.adaptors.Zigbee2mqtt.Update(model); err != nil {
		return
	}

	result, err = z.adaptors.Zigbee2mqtt.GetById(model.Id)
	bridge.UpdateModel(result)

	return
}

func (z *Zigbee2mqtt) DeleteBridge(bridgeId int64) (err error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	var bridge *Bridge
	if bridge, err = z.unsafeGetBridge(bridgeId); err == nil {
		bridge.Stop(context.Background())
		delete(z.bridges, bridgeId)
	} else {
		return
	}

	err = z.adaptors.Zigbee2mqtt.Delete(bridgeId)

	return
}

func (z *Zigbee2mqtt) ResetBridge(bridgeId int64) (err error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	var bridge *Bridge
	if bridge, err = z.unsafeGetBridge(bridgeId); err == nil {
		bridge.ConfigReset()
	}
	return
}

func (z *Zigbee2mqtt) BridgeDeviceBan(bridgeId int64, friendlyName string) (err error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	var bridge *Bridge
	if bridge, err = z.unsafeGetBridge(bridgeId); err == nil {
		bridge.Ban(friendlyName)
	}
	return
}

func (z *Zigbee2mqtt) BridgeDeviceWhitelist(bridgeId int64, friendlyName string) (err error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	var bridge *Bridge
	if bridge, err = z.unsafeGetBridge(bridgeId); err == nil {
		bridge.Whitelist(friendlyName)
	}
	return
}

func (z *Zigbee2mqtt) BridgeNetworkmap(bridgeId int64) (networkmap string, err error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	var bridge *Bridge
	if bridge, err = z.unsafeGetBridge(bridgeId); err == nil {
		networkmap = bridge.Networkmap()
	}
	return
}

func (z *Zigbee2mqtt) BridgeUpdateNetworkmap(bridgeId int64) (err error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	var bridge *Bridge
	if bridge, err = z.unsafeGetBridge(bridgeId); err == nil {
		bridge.UpdateNetworkmap()
	}
	return
}

func (z *Zigbee2mqtt) unsafeGetBridge(bridgeId int64) (bridge *Bridge, err error) {
	var ok bool
	if bridge, ok = z.bridges[bridgeId]; !ok {
		err = adaptors.ErrRecordNotFound
	}
	return
}

func (z *Zigbee2mqtt) GetTopicByDevice(model *m.Zigbee2mqttDevice) (topic string, err error) {

	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	br, ok := z.bridges[model.Zigbee2mqttId]
	if !ok {
		err = adaptors.ErrRecordNotFound
		return
	}

	topic = br.GetDeviceTopic(model.Id)

	return
}

func (z *Zigbee2mqtt) DeviceRename(friendlyName, name string) (err error) {
	z.bridgesLock.Lock()
	defer z.bridgesLock.Unlock()

	for _, bridge := range z.bridges {
		_ = bridge.RenameDevice(friendlyName, name)
	}

	return
}
