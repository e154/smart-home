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

type Flow struct {
	Total    int64 `json:"total"`
	Disabled int64 `json:"disabled"`
}

type FlowManager struct {
	publisher IPublisher
	total     metrics.Counter
	disabled  metrics.Counter
}

func NewFlowManager(publisher IPublisher) *FlowManager {
	return &FlowManager{
		publisher: publisher,
		total:     metrics.NewCounter(),
		disabled:  metrics.NewCounter(),
	}
}

func (d *FlowManager) update(t interface{}) {
	switch v := t.(type) {
	case FlowAdd:
		d.total.Inc(v.TotalNum)
		d.disabled.Inc(v.DisabledNum)
	case FlowDelete:
		d.total.Dec(v.TotalNum)
		d.disabled.Dec(v.DisabledNum)

	default:
		return
	}

	d.broadcast()
}

func (d *FlowManager) Snapshot() Flow {

	return Flow{
		Total:    d.total.Count(),
		Disabled: d.disabled.Count(),
	}
}

func (d *FlowManager) broadcast() {
	go d.publisher.Broadcast("flow")
}

type FlowAdd struct {
	TotalNum    int64
	DisabledNum int64
}

type FlowDelete struct {
	TotalNum    int64
	DisabledNum int64
}
