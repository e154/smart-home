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
	"github.com/e154/smart-home/adaptors"
	"github.com/rcrowley/go-metrics"
)

// Flow ...
type Flow struct {
	Total    int64 `json:"total"`
	Disabled int64 `json:"disabled"`
}

// FlowManager ...
type FlowManager struct {
	publisher IPublisher
	total     metrics.Counter
	enabled   metrics.Counter
}

// NewFlowManager ...
func NewFlowManager(publisher IPublisher,
	adaptors *adaptors.Adaptors) (flow *FlowManager) {

	flow = &FlowManager{
		publisher: publisher,
		total:     metrics.NewCounter(),
		enabled:   metrics.NewCounter(),
	}

	if _, total, err := adaptors.Flow.List(999, 0, "", ""); err == nil {
		flow.total.Inc(total)
	}

	return
}

func (d *FlowManager) update(t interface{}) {
	switch v := t.(type) {
	case FlowAdd:
		d.total.Inc(v.TotalNum)
		d.enabled.Inc(v.EnabledNum)
	case FlowDelete:
		d.total.Dec(v.TotalNum)
		d.enabled.Dec(v.EnabledNum)

	default:
		return
	}

	d.broadcast()
}

// Snapshot ...
func (d *FlowManager) Snapshot() Flow {

	return Flow{
		Total:    d.total.Count(),
		Disabled: d.total.Count() - d.enabled.Count(),
	}
}

func (d *FlowManager) broadcast() {
	go d.publisher.Broadcast("flow")
}

// FlowAdd ...
type FlowAdd struct {
	TotalNum   int64
	EnabledNum int64
}

// FlowDelete ...
type FlowDelete struct {
	TotalNum   int64
	EnabledNum int64
}
