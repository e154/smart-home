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

package adaptors

type Adaptors struct {
	Script            ScriptRepo
	Tag               TagRepo
	Role              RoleRepo
	Permission        PermissionRepo
	User              UserRepo
	UserMeta          UserMetaRepo
	UserDevice        UserDeviceRepo
	Image             ImageRepo
	Variable          VariableRepo
	Entity            EntityRepo
	EntityState       EntityStateRepo
	EntityAction      EntityActionRepo
	EntityStorage     EntityStorageRepo
	Log               LogRepo
	Template          TemplateRepo
	Message           MessageRepo
	MessageDelivery   MessageDeliveryRepo
	Zigbee2mqtt       Zigbee2mqttRepo
	Zigbee2mqttDevice Zigbee2mqttDeviceRepo
	AlexaSkill        AlexaSkillRepo
	AlexaIntent       AlexaIntentRepo
	Metric            MetricRepo
	MetricBucket      MetricBucketRepo
	Area              AreaRepo
	Action            ActionRepo
	Condition         ConditionRepo
	Trigger           TriggerRepo
	Task              TaskRepo
	RunHistory        RunHistoryRepo
	Plugin            PluginRepo
	TelegramChat      TelegramChatRepo
	Dashboard         DashboardRepo
	DashboardTab      DashboardTabRepo
	DashboardCard     IDashboardCard
	DashboardCardItem DashboardCardItemRepo
	ScriptVersion     ScriptVersionRepo
	Automation        AutomationRepo
	Transaction       TransactionManger
}
