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

package uptime

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	// EntitySensor ...
	EntitySensor = string("uptime")

	// Name ...
	Name = "uptime"

	Version = "0.0.1"

	AttrUptimeTotal         = "uptime_total"
	AttrAppStarted          = "app_started"
	AttrFirstStart          = "first_start"
	AttrLastShutdown        = "last_shutdown"
	AttrLastShutdownCorrect = "last_shutdown_correct"
	AttrLastStart           = "last_start"
	AttrUptime              = "uptime"
	AttrDowntime            = "downtime"
	AttrUptimePercent       = "uptime_percent"
	AttrDowntimePercent     = "downtime_percent"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrUptimeTotal: {
			Name: AttrUptimeTotal,
			Type: common.AttributeInt,
		},
		AttrAppStarted: {
			Name: AttrAppStarted,
			Type: common.AttributeTime,
		},
		AttrFirstStart: {
			Name: AttrFirstStart,
			Type: common.AttributeTime,
		},
		AttrLastShutdown: {
			Name: AttrLastShutdown,
			Type: common.AttributeTime,
		},
		AttrLastShutdownCorrect: {
			Name: AttrLastShutdownCorrect,
			Type: common.AttributeBool,
		},
		AttrLastStart: {
			Name: AttrLastStart,
			Type: common.AttributeTime,
		},
		AttrUptime: {
			Name: AttrUptime,
			Type: common.AttributeInt,
		},
		AttrDowntime: {
			Name: AttrDowntime,
			Type: common.AttributeInt,
		},
		AttrUptimePercent: {
			Name: AttrUptimePercent,
			Type: common.AttributeFloat,
		},
		AttrDowntimePercent: {
			Name: AttrDowntimePercent,
			Type: common.AttributeFloat,
		},
	}
}
