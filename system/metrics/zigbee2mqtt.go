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

package metrics

import (
	"github.com/rcrowley/go-metrics"
)

type Zigbee2Mqtt struct {
	Total    int64 `json:"total"`
	Disabled int64 `json:"disabled"`
}

type Zigbee2MqttManager struct {
	publisher IPublisher
	total     metrics.Counter
	disabled  metrics.Counter
}

func NewZigbee2MqttManager(publisher IPublisher) *Zigbee2MqttManager {
	return &Zigbee2MqttManager{
		publisher: publisher,
		total:     metrics.NewCounter(),
		disabled:  metrics.NewCounter(),
	}
}

func (d *Zigbee2MqttManager) update(t interface{}) {
	switch v := t.(type) {
	case Zigbee2MqttAdd:
		d.total.Inc(v.TotalNum)
	case Zigbee2MqttDelete:
		d.total.Dec(v.TotalNum)
		d.disabled.Dec(v.DisabledNum)
	default:
		return
	}

	d.broadcast()
}

func (d *Zigbee2MqttManager) Snapshot() Zigbee2Mqtt {
	return Zigbee2Mqtt{
		Total:    d.total.Count(),
		Disabled: d.disabled.Count(),
	}
}

func (d *Zigbee2MqttManager) broadcast() {
	go d.publisher.Broadcast("zigbee2mqtt")
}

type Zigbee2MqttAdd struct {
	TotalNum    int64
	DisabledNum int64
}

type Zigbee2MqttDelete struct {
	TotalNum    int64
	DisabledNum int64
}
