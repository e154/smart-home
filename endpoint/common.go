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

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/cache"
	"github.com/e154/smart-home/system/jwt_manager"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/system/zigbee2mqtt"
)

// CommonEndpoint ...
type CommonEndpoint struct {
	adaptors      *adaptors.Adaptors
	accessList    access_list.AccessListService
	scriptService scripts.ScriptService
	zigbee2mqtt   zigbee2mqtt.Zigbee2mqtt
	eventBus      bus.Bus
	supervisor    supervisor.Supervisor
	mqtt          mqtt.MqttServ
	jwtManager    jwt_manager.JwtManager
	validation    *validation.Validate
	appConfig     *m.AppConfig
	automation    automation.Automation
	cache         cache.Cache
}

// NewCommonEndpoint ...
func NewCommonEndpoint(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	scriptService scripts.ScriptService,
	zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
	eventBus bus.Bus,
	supervisor supervisor.Supervisor,
	mqtt mqtt.MqttServ,
	jwtManager jwt_manager.JwtManager,
	validation *validation.Validate,
	appConfig *m.AppConfig,
	automation automation.Automation,
) *CommonEndpoint {
	cache, _ := cache.NewCache("memory", `{"interval":60}`)
	return &CommonEndpoint{
		adaptors:      adaptors,
		accessList:    accessList,
		scriptService: scriptService,
		zigbee2mqtt:   zigbee2mqtt,
		eventBus:      eventBus,
		supervisor:    supervisor,
		mqtt:          mqtt,
		jwtManager:    jwtManager,
		validation:    validation,
		appConfig:     appConfig,
		automation:    automation,
		cache:         cache,
	}
}

func (c *CommonEndpoint) checkSuperUser(ctx context.Context) (decline bool) {
	root, _ := ctx.Value("root").(bool)
	//log.Debugf("root: %t, %t", root, ok)

	return !root
}
