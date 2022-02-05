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

package endpoint

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/jwt_manager"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/system/zigbee2mqtt"
)

var (
	log = common.MustGetLogger("endpoint")
)

// Endpoint ...
type Endpoint struct {
	AlexaSkill      *AlexaSkillEndpoint
	Auth            *AuthEndpoint
	Image           *ImageEndpoint
	Log             *LogEndpoint
	Role            *RoleEndpoint
	Script          *ScriptEndpoint
	User            *UserEndpoint
	Template        *TemplateEndpoint
	Notify          *NotifyEndpoint
	MessageDelivery *MessageDeliveryEndpoint
	Version         *VersionEndpoint
	Zigbee2mqtt     *Zigbee2mqttEndpoint
	Entity          *EntityEndpoint
	DeveloperTools  *DeveloperToolsEndpoint
	Mqtt            *MqttEndpoint
	Plugin          *PluginEndpoint
	PluginActor     *PluginActorEndpoint
	Task            *TaskEndpoint
	Area            *AreaEndpoint
	Interact        *InteractEndpoint
}

// NewEndpoint ...
func NewEndpoint(adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	accessList access_list.AccessListService,
	zigbee2mqtt zigbee2mqtt.Zigbee2mqtt,
	entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus,
	pluginManager common.PluginManager,
	mqtt mqtt.MqttServ,
	jwtManager jwt_manager.JwtManager,
	validation *validation.Validate) *Endpoint {
	common := NewCommonEndpoint(adaptors, accessList, scriptService, zigbee2mqtt, eventBus, pluginManager, entityManager, mqtt, jwtManager, validation)
	return &Endpoint{
		AlexaSkill:      NewAlexaSkillEndpoint(common),
		Auth:            NewAuthEndpoint(common),
		Image:           NewImageEndpoint(common),
		Log:             NewLogEndpoint(common),
		Role:            NewRoleEndpoint(common),
		Script:          NewScriptEndpoint(common),
		User:            NewUserEndpoint(common),
		Template:        NewTemplateEndpoint(common),
		Notify:          NewNotifyEndpoint(common),
		MessageDelivery: NewMessageDeliveryEndpoint(common),
		Version:         NewVersionEndpoint(common),
		Zigbee2mqtt:     NewZigbee2mqttEndpoint(common),
		Entity:          NewEntityEndpoint(common),
		DeveloperTools:  NewDeveloperToolsEndpoint(common),
		Mqtt:            NewMqttEndpoint(common),
		Plugin:          NewPluginEndpoint(common),
		PluginActor:     NewPluginActorEndpoint(common),
		Task:            NewTaskEndpoint(common),
		Area:            NewAreaEndpoint(common),
		Interact:        NewInteractEndpoint(common),
	}
}
