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
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/system/backup"
	"github.com/e154/smart-home/system/stream"
)

var (
	log = logger.MustGetLogger("endpoint")
)

// Endpoint ...
type Endpoint struct {
	AlexaSkill        *AlexaSkillEndpoint
	Auth              *AuthEndpoint
	Image             *ImageEndpoint
	Log               *LogEndpoint
	Role              *RoleEndpoint
	Script            *ScriptEndpoint
	Tag               *TagEndpoint
	User              *UserEndpoint
	Template          *TemplateEndpoint
	Notify            *NotifyEndpoint
	MessageDelivery   *MessageDeliveryEndpoint
	Version           *VersionEndpoint
	Zigbee2mqtt       *Zigbee2mqttEndpoint
	Entity            *EntityEndpoint
	DeveloperTools    *DeveloperToolsEndpoint
	Mqtt              *MqttEndpoint
	Plugin            *PluginEndpoint
	PluginActor       *PluginActorEndpoint
	Action            *ActionEndpoint
	Condition         *ConditionEndpoint
	Trigger           *TriggerEndpoint
	Task              *TaskEndpoint
	Area              *AreaEndpoint
	Interact          *InteractEndpoint
	Dashboard         *DashboardEndpoint
	DashboardTab      *DashboardTabEndpoint
	DashboardCard     *DashboardCardEndpoint
	DashboardCardItem *DashboardCardItemEndpoint
	Variable          *VariableEndpoint
	EntityStorage     *EntityStorageEndpoint
	Metric            *MetricEndpoint
	Backup            *BackupEndpoint
	Stream            *StreamEndpoint
	Webdav            *WebdavEndpoint
	Webhook           *WebhookEndpoint
	Automation        *AutomationEndpoint
}

// NewEndpoint ...
func NewEndpoint(backup *backup.Backup, stream *stream.Stream, common *CommonEndpoint) *Endpoint {
	return &Endpoint{
		AlexaSkill:        NewAlexaSkillEndpoint(common),
		Auth:              NewAuthEndpoint(common),
		Image:             NewImageEndpoint(common),
		Log:               NewLogEndpoint(common),
		Role:              NewRoleEndpoint(common),
		Script:            NewScriptEndpoint(common),
		Tag:               NewTagEndpoint(common),
		User:              NewUserEndpoint(common),
		Template:          NewTemplateEndpoint(common),
		Notify:            NewNotifyEndpoint(common),
		MessageDelivery:   NewMessageDeliveryEndpoint(common),
		Version:           NewVersionEndpoint(common),
		Zigbee2mqtt:       NewZigbee2mqttEndpoint(common),
		Entity:            NewEntityEndpoint(common),
		DeveloperTools:    NewDeveloperToolsEndpoint(common),
		Mqtt:              NewMqttEndpoint(common),
		Plugin:            NewPluginEndpoint(common),
		PluginActor:       NewPluginActorEndpoint(common),
		Action:            NewActionEndpoint(common),
		Condition:         NewConditionEndpoint(common),
		Trigger:           NewTriggerEndpoint(common),
		Task:              NewTaskEndpoint(common),
		Area:              NewAreaEndpoint(common),
		Interact:          NewInteractEndpoint(common),
		Dashboard:         NewDashboardEndpoint(common),
		DashboardTab:      NewDashboardTabEndpoint(common),
		DashboardCard:     NewDashboardCardEndpoint(common),
		DashboardCardItem: NewDashboardCardItemEndpoint(common),
		Variable:          NewVariableEndpoint(common),
		EntityStorage:     NewEntityStorageEndpoint(common),
		Metric:            NewMetricEndpoint(common),
		Backup:            NewBackupEndpoint(common, backup),
		Stream:            NewStreamEndpoint(common, stream),
		Webdav:            NewWebdavEndpoint(common),
		Webhook:           NewWebhookEndpoint(common),
		Automation:        NewAutomationEndpoint(common),
	}
}
