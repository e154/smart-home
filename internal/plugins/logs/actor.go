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

package logs

import (
	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/rcrowley/go-metrics"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	cores         int64
	model         string
	ErrTotal      metrics.Counter
	ErrToday      metrics.Counter
	ErrYesterday  metrics.Counter
	WarnTotal     metrics.Counter
	WarnToday     metrics.Counter
	WarnYesterday metrics.Counter
}

// NewActor ...
func NewActor(entity *m.Entity,
	service plugins.Service) *Actor {

	actor := &Actor{
		BaseActor:     supervisor.NewBaseActor(entity, service),
		ErrTotal:      metrics.NewCounter(),
		ErrToday:      metrics.NewCounter(),
		ErrYesterday:  metrics.NewCounter(),
		WarnTotal:     metrics.NewCounter(),
		WarnToday:     metrics.NewCounter(),
		WarnYesterday: metrics.NewCounter(),
	}

	if entity != nil {
		actor.Metric = entity.Metrics
		attrs := entity.Attributes
		actor.ErrTotal.Inc(attrs[AttrErrTotal].Int64())
		actor.ErrToday.Inc(attrs[AttrErrToday].Int64())
		actor.ErrYesterday.Inc(attrs[AttrErrYesterday].Int64())
		actor.WarnTotal.Inc(attrs[AttrWarnTotal].Int64())
		actor.WarnToday.Inc(attrs[AttrWarnToday].Int64())
		actor.WarnYesterday.Inc(attrs[AttrWarnYesterday].Int64())
	}
	return actor
}

func (e *Actor) Destroy() {

}

func (e *Actor) Spawn() {
	go e.selfUpdate()
}

func (e *Actor) selfUpdate() {

	e.AttrMu.Lock()
	e.Attrs[AttrErrTotal].Value = e.ErrTotal.Count()
	e.Attrs[AttrErrToday].Value = e.ErrToday.Count()
	e.Attrs[AttrErrYesterday].Value = e.ErrYesterday.Count()
	e.Attrs[AttrWarnTotal].Value = e.WarnTotal.Count()
	e.Attrs[AttrWarnToday].Value = e.WarnToday.Count()
	e.Attrs[AttrWarnYesterday].Value = e.WarnYesterday.Count()
	e.AttrMu.Unlock()

	e.SaveState(false, true)
}

func (e *Actor) LogsHook(level common.LogLevel) {

	switch level {
	case common.LogLevelError:
		e.ErrTotal.Inc(1)
		e.ErrToday.Inc(1)
	case common.LogLevelWarning:
		e.WarnTotal.Inc(1)
		e.WarnToday.Inc(1)
	//case common.LogLevelInfo:
	//case common.LogLevelDebug:
	default:
		return
	}
	e.selfUpdate()
}

func (e *Actor) UpdateDay() {
	e.ErrYesterday.Clear()
	e.ErrYesterday.Inc(e.ErrToday.Count())
	e.WarnYesterday.Clear()
	e.WarnYesterday.Inc(e.WarnToday.Count())
	e.ErrToday.Clear()
	e.WarnToday.Clear()

	e.selfUpdate()
}
