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

package endpoint

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/notify"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
)

// CommonEndpoint ...
type CommonEndpoint struct {
	adaptors      *adaptors.Adaptors
	accessList    access_list.AccessListService
	scriptService scripts.ScriptService
	notify        notify.Notify
	zigbee2mqtt   zigbee2mqtt.Zigbee2mqtt
	mqtt          mqtt.MqttServ
}

// NewCommonEndpoint ...
func NewCommonEndpoint(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	scriptService scripts.ScriptService,
	notify notify.Notify,
	zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
	mqtt mqtt.MqttServ,
) *CommonEndpoint {
	return &CommonEndpoint{
		adaptors:      adaptors,
		accessList:    accessList,
		scriptService: scriptService,
		notify:        notify,
		zigbee2mqtt:   zigbee2mqtt,
		mqtt:          mqtt,
	}
}
