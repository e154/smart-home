// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

package time

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	FunctionName   = "automationTriggerTime"
	Name           = "time"
	Version        = "0.0.1"
	AttrCron       = "cron"
	AttrSystemInfo = "SystemInfo"
)

func NewTriggerParams() m.TriggerParams {
	return m.TriggerParams{
		Script:   true,
		Required: []string{AttrCron},
		Attributes: m.Attributes{
			AttrSystemInfo: {
				Name: AttrSystemInfo,
				Type: common.AttributeNotice,
			},
			AttrCron: {
				Name: AttrCron,
				Type: common.AttributeString,
			},
		},
	}
}
