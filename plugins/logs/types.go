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

package logs

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	// EntityLogs ...
	EntityLogs = string("logs")

	AttrErrTotal      = "err_total"
	AttrErrToday      = "err_today"
	AttrErrYesterday  = "err_yesterday"
	AttrWarnTotal     = "warn_total"
	AttrWarnToday     = "warn_today"
	AttrWarnYesterday = "warn_yesterday"

	// Name ...
	Name = "logs"

	// EntityType ...
	EntityType = "logs"

	Version = "0.0.1"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrErrTotal: {
			Name: AttrErrTotal,
			Type: common.AttributeInt,
		},
		AttrErrToday: {
			Name: AttrErrToday,
			Type: common.AttributeInt,
		},
		AttrErrYesterday: {
			Name: AttrErrYesterday,
			Type: common.AttributeInt,
		},
		AttrWarnTotal: {
			Name: AttrWarnTotal,
			Type: common.AttributeInt,
		},
		AttrWarnToday: {
			Name: AttrWarnToday,
			Type: common.AttributeInt,
		},
		AttrWarnYesterday: {
			Name: AttrWarnYesterday,
			Type: common.AttributeInt,
		},
	}
}
