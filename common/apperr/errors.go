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

package apperr

import (
	"errors"
)

var (
	ErrInternal        = errors.New("internal error")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotAllowed      = errors.New("not allowed")
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInvalidRequest  = errors.New("invalid request")
	ErrUnauthorized    = errors.New("unauthorized")
)

var (
	ErrUnimplemented              = errors.New("unimplemented")
	ErrPassNotValid               = errors.New("password not valid")
	ErrAccountIsBlocked           = errors.New("account is blocked")
	ErrTokenIsDeprecated          = errors.New("token is deprecated")
	ErrBadLoginOrPassword         = errors.New("bad login or password")
	ErrMqttServerNoWorked         = errors.New("mqtt server not worked")
	ErrBadRequestParams           = errors.New("bad request params")
	ErrBadActorSettingsParameters = errors.New("bad actor settings parameters")
	ErrTimeout                    = errors.New("timeout")
	ErrProviderIsEmpty            = errors.New("provider is empty")
	ErrBadSettings                = errors.New("there are no settings or the parameters are incorrectly set")
)

var (
	ErrDashboardImport   = New("failed to import dashboard", ErrInternal)
	ErrDashboardAdd      = New("failed to add dashboard", ErrInternal)
	ErrDashboardGet      = New("failed to get dashboard", ErrInternal)
	ErrDashboardList     = New("failed to list dashboard", ErrInternal)
	ErrDashboardNotFound = New("dashboard is not found", ErrNotFound)
	ErrDashboardUpdate   = New("failed to update dashboard", ErrInternal)
	ErrDashboardDelete   = New("failed to delete dashboard", ErrInternal)
	ErrDashboardSearch   = New("failed to search dashboard", ErrInternal)

	ErrDashboardCardAdd      = New("failed to add dashboard card", ErrInternal)
	ErrDashboardCardGet      = New("failed to get dashboard card", ErrInternal)
	ErrDashboardCardList     = New("failed to list dashboard card", ErrInternal)
	ErrDashboardCardNotFound = New("dashboard card is not found", ErrNotFound)
	ErrDashboardCardUpdate   = New("failed to update dashboard card", ErrInternal)
	ErrDashboardCardDelete   = New("failed to delete dashboard card", ErrInternal)

	ErrDashboardCardItemAdd      = New("failed to add dashboard card item", ErrInternal)
	ErrDashboardCardItemGet      = New("failed to get dashboard card item", ErrInternal)
	ErrDashboardCardItemList     = New("failed to list dashboard card item", ErrInternal)
	ErrDashboardCardItemNotFound = New("dashboard card item is not found", ErrNotFound)
	ErrDashboardCardItemUpdate   = New("failed to update dashboard card item", ErrInternal)
	ErrDashboardCardItemDelete   = New("failed to delete dashboard card item", ErrInternal)

	ErrDashboardTabAdd      = New("failed to add dashboard tab", ErrInternal)
	ErrDashboardTabGet      = New("failed to get dashboard tab", ErrInternal)
	ErrDashboardTabList     = New("failed to list dashboard tab", ErrInternal)
	ErrDashboardTabNotFound = New("dashboard tab is not found", ErrNotFound)
	ErrDashboardTabUpdate   = New("failed to update dashboard tab", ErrInternal)
	ErrDashboardTabDelete   = New("failed to delete dashboard tab", ErrInternal)

	ErrActionAdd      = New("failed to add action", ErrInternal)
	ErrActionGet      = New("failed to get action", ErrInternal)
	ErrActionUpdate   = New("failed to update action", ErrInternal)
	ErrActionList     = New("failed to list action", ErrInternal)
	ErrActionNotFound = New("action is not found", ErrNotFound)
	ErrActionDelete   = New("failed to delete action", ErrInternal)
	ErrActionSearch   = New("failed to search action", ErrInternal)

	ErrEntityAdd           = New("failed to add entity", ErrInternal)
	ErrEntityGet           = New("failed to get entity", ErrInternal)
	ErrEntityList          = New("failed to list entity", ErrInternal)
	ErrEntityNotFound      = New("entity is not found", ErrNotFound)
	ErrEntityUpdate        = New("failed to update entity", ErrInternal)
	ErrEntityDelete        = New("failed to delete entity", ErrInternal)
	ErrEntitySerch         = New("failed to search entity", ErrInternal)
	ErrEntityAppendMetric  = New("entity append metric", ErrInternal)
	ErrEntityDeleteMetric  = New("delete metric failed", ErrInternal)
	ErrEntityReplaceMetric = New("replace metric failed", ErrInternal)
	ErrEntityAppendScript  = New("append script failed", ErrInternal)
	ErrEntityDeleteScript  = New("delete script failed", ErrInternal)
	ErrEntityReplaceScript = New("replace script failed", ErrInternal)

	ErrAlexaIntentAdd      = New("failed to add intent", ErrInternal)
	ErrAlexaIntentUpdate   = New("failed to update intent", ErrInternal)
	ErrAlexaIntentGet      = New("failed to get intent", ErrInternal)
	ErrAlexaIntentDelete   = New("failed to delete intent", ErrInternal)
	ErrAlexaIntentNotFound = New("intent is not found", ErrNotFound)

	ErrAlexaSkillAdd      = New("failed to add skill", ErrInternal)
	ErrAlexaSkillGet      = New("failed to get skill", ErrInternal)
	ErrAlexaSkillUpdate   = New("failed to update skill", ErrInternal)
	ErrAlexaSkillList     = New("failed to list skill", ErrInternal)
	ErrAlexaSkillNotFound = New("skill is not found", ErrNotFound)
	ErrAlexaSkillDelete   = New("failed to delete skill", ErrInternal)

	ErrAreaAdd      = New("failed to add area", ErrInternal)
	ErrAreaGet      = New("failed to get area", ErrInternal)
	ErrAreaUpdate   = New("failed to update area", ErrInternal)
	ErrAreaList     = New("failed to list area", ErrInternal)
	ErrAreaNotFound = New("area is not found", ErrNotFound)
	ErrAreaDelete   = New("failed to delete area", ErrInternal)
	ErrAreaClean    = New("failed to clean area", ErrInternal)

	ErrConditionAdd      = New("failed to add condition", ErrInternal)
	ErrConditionGet      = New("failed to get condition", ErrInternal)
	ErrConditionUpdate   = New("failed to update condition", ErrInternal)
	ErrConditionList     = New("failed to list condition", ErrInternal)
	ErrConditionNotFound = New("condition is not found", ErrNotFound)
	ErrConditionDelete   = New("failed to delete condition", ErrInternal)
	ErrConditionSearch   = New("failed to search condition", ErrInternal)

	ErrEntityActionAdd      = New("failed to add action", ErrInternal)
	ErrEntityActionGet      = New("failed to get action", ErrInternal)
	ErrEntityActionUpdate   = New("failed to update action", ErrInternal)
	ErrEntityActionList     = New("failed to list action", ErrInternal)
	ErrEntityActionNotFound = New("action is not found", ErrNotFound)
	ErrEntityActionDelete   = New("failed to delete action", ErrInternal)

	ErrEntityStateAdd      = New("failed to add state", ErrInternal)
	ErrEntityStateGet      = New("failed to get state", ErrInternal)
	ErrEntityStateUpdate   = New("failed to update state", ErrInternal)
	ErrEntityStateList     = New("failed to list state", ErrInternal)
	ErrEntityStateNotFound = New("state is not found", ErrNotFound)
	ErrEntityStateDelete   = New("failed to delete state", ErrInternal)

	ErrEntityStorageAdd    = New("failed to add storage", ErrInternal)
	ErrEntityStorageGet    = New("failed to get storage", ErrInternal)
	ErrEntityStorageList   = New("failed to list storage", ErrInternal)
	ErrEntityStorageDelete = New("failed to delete storage", ErrInternal)

	ErrImageAdd      = New("failed to add image", ErrInternal)
	ErrImageGet      = New("failed to get image", ErrInternal)
	ErrImageUpdate   = New("failed to update image", ErrInternal)
	ErrImageList     = New("failed to list image", ErrInternal)
	ErrImageNotFound = New("image is not found", ErrNotFound)
	ErrImageDelete   = New("failed to delete image", ErrInternal)

	ErrLogAdd      = New("failed to add log", ErrInternal)
	ErrLogGet      = New("failed to get log", ErrInternal)
	ErrLogList     = New("failed to list log", ErrInternal)
	ErrLogNotFound = New("log is not found", ErrNotFound)
	ErrLogDelete   = New("failed to delete log", ErrNotFound)

	ErrMessageAdd              = New("failed to add message", ErrInternal)
	ErrMessageDeliveryAdd      = New("failed to add message delivery", ErrInternal)
	ErrMessageDeliveryList     = New("failed to list message delivery", ErrInternal)
	ErrMessageDeliveryUpdate   = New("failed to update message delivery", ErrInternal)
	ErrMessageDeliveryDelete   = New("failed to delete message delivery", ErrInternal)
	ErrMessageDeliveryGet      = New("failed to get message delivery", ErrInternal)
	ErrMessageDeliveryNotFound = New("message delivery is not found", ErrNotFound)

	ErrMetricAdd      = New("failed to add metric", ErrInternal)
	ErrMetricGet      = New("failed to get metric", ErrInternal)
	ErrMetricUpdate   = New("failed to update metric", ErrInternal)
	ErrMetricList     = New("failed to list metric", ErrInternal)
	ErrMetricNotFound = New("metric is not found", ErrNotFound)
	ErrMetricDelete   = New("failed to delete metric", ErrInternal)
	ErrMetricSearch   = New("failed to search metric", ErrInternal)

	ErrMetricBucketAdd    = New("failed to add metric backet", ErrInternal)
	ErrMetricBucketGet    = New("failed to get metric backet", ErrInternal)
	ErrMetricBucketDelete = New("failed to delete metric backet", ErrInternal)

	ErrPermissionAdd    = New("failed to add permission", ErrInternal)
	ErrPermissionGet    = New("failed to get permission", ErrInternal)
	ErrPermissionDelete = New("failed to delete permission", ErrInternal)

	ErrPluginAdd      = New("failed to add plugin", ErrInternal)
	ErrPluginGet      = New("failed to get plugin", ErrInternal)
	ErrPluginUpdate   = New("failed to update plugin", ErrInternal)
	ErrPluginList     = New("failed to list plugin", ErrInternal)
	ErrPluginNotFound = New("plugin is not found", ErrNotFound)
	ErrPluginDelete   = New("failed to delete plugin", ErrInternal)
	ErrPluginSearch   = New("failed to search plugin", ErrInternal)

	ErrRoleAdd      = New("failed to add role", ErrInternal)
	ErrRoleGet      = New("failed to get role", ErrInternal)
	ErrRoleUpdate   = New("failed to update role", ErrInternal)
	ErrRoleList     = New("failed to list role", ErrInternal)
	ErrRoleNotFound = New("role is not found", ErrNotFound)
	ErrRoleDelete   = New("failed to delete role", ErrInternal)
	ErrRoleSearch   = New("failed to search role", ErrInternal)

	ErrRunStoryAdd    = New("failed to add run story", ErrInternal)
	ErrRunStoryUpdate = New("failed to update run story", ErrInternal)
	ErrRunStoryList   = New("failed to list run story", ErrInternal)

	ErrScriptAdd      = New("failed to add script", ErrInternal)
	ErrScriptGet      = New("failed to get script", ErrInternal)
	ErrScriptUpdate   = New("failed to update script", ErrInternal)
	ErrScriptList     = New("failed to list script", ErrInternal)
	ErrScriptNotFound = New("script is not found", ErrNotFound)
	ErrScriptDelete   = New("failed to delete script", ErrInternal)
	ErrScriptSearch   = New("failed to search script", ErrInternal)
	ErrScriptStat     = New("failed to get script statistic", ErrInternal)

	ErrTaskAdd             = New("failed to add Task", ErrInternal)
	ErrTaskGet             = New("failed to get Task", ErrInternal)
	ErrTaskUpdate          = New("failed to update Task", ErrInternal)
	ErrTaskList            = New("failed to list Task", ErrInternal)
	ErrTaskNotFound        = New("Task is not found", ErrNotFound)
	ErrTaskDelete          = New("failed to delete Task", ErrInternal)
	ErrTaskSearch          = New("failed to search Task", ErrInternal)
	ErrTaskAppendTrigger   = New("task append trigger failed", ErrInternal)
	ErrTaskDeleteTrigger   = New("task delete trigger failed", ErrInternal)
	ErrTaskAppendCondition = New("task append condition failed", ErrInternal)
	ErrTaskDeleteCondition = New("task delete condition failed", ErrInternal)
	ErrTaskAppendAction    = New("task append action failed", ErrInternal)
	ErrTaskDeleteAction    = New("task delete action failed", ErrInternal)

	ErrChatAdd    = New("failed to add chat", ErrInternal)
	ErrChatList   = New("failed to list chat", ErrInternal)
	ErrChatDelete = New("failed to delete chat", ErrInternal)

	ErrTemplateAdd      = New("failed to add template", ErrInternal)
	ErrTemplateGet      = New("failed to get template", ErrInternal)
	ErrTemplateUpdate   = New("failed to update template", ErrInternal)
	ErrTemplateList     = New("failed to list template", ErrInternal)
	ErrTemplateNotFound = New("template is not found", ErrNotFound)
	ErrTemplateDelete   = New("failed to delete template", ErrInternal)
	ErrTemplateSearch   = New("failed to search template", ErrInternal)

	ErrTriggerAdd      = New("failed to add trigger", ErrInternal)
	ErrTriggerGet      = New("failed to get trigger", ErrInternal)
	ErrTriggerUpdate   = New("failed to update trigger", ErrInternal)
	ErrTriggerList     = New("failed to list trigger", ErrInternal)
	ErrTriggerNotFound = New("trigger is not found", ErrNotFound)
	ErrTriggerDelete   = New("failed to delete trigger", ErrInternal)
	ErrTriggerSearch   = New("failed to search trigger", ErrInternal)

	ErrUserAdd      = New("failed to add user", ErrInternal)
	ErrUserMetaAdd  = New("failed to add user meta", ErrInternal)
	ErrUserGet      = New("failed to get user", ErrInternal)
	ErrUserUpdate   = New("failed to update user", ErrInternal)
	ErrUserList     = New("failed to list user", ErrInternal)
	ErrUserNotFound = New("user is not found", ErrNotFound)
	ErrUserDelete   = New("failed to delete user", ErrInternal)

	ErrVariableAdd      = New("failed to add variable", ErrInternal)
	ErrVariableGet      = New("failed to get variable", ErrInternal)
	ErrVariableUpdate   = New("failed to update variable", ErrInternal)
	ErrVariableList     = New("failed to list variable", ErrInternal)
	ErrVariableNotFound = New("variable is not found", ErrNotFound)
	ErrVariableDelete   = New("failed to delete variable", ErrInternal)

	ErrZigbee2mqttAdd      = New("failed to add zigbee2mqtt", ErrInternal)
	ErrZigbee2mqttGet      = New("failed to get zigbee2mqtt", ErrInternal)
	ErrZigbee2mqttUpdate   = New("failed to update zigbee2mqtt", ErrInternal)
	ErrZigbee2mqttList     = New("failed to list zigbee2mqtt", ErrInternal)
	ErrZigbee2mqttNotFound = New("zigbee2mqtt is not found", ErrNotFound)
	ErrZigbee2mqttDelete   = New("failed to delete zigbee2mqtt", ErrInternal)

	ErrZigbeeDeviceAdd      = New("failed to add device", ErrInternal)
	ErrZigbeeDeviceGet      = New("failed to get device", ErrInternal)
	ErrZigbeeDeviceUpdate   = New("failed to update device", ErrInternal)
	ErrZigbeeDeviceList     = New("failed to list device", ErrInternal)
	ErrZigbeeDeviceNotFound = New("device is not found", ErrNotFound)
	ErrZigbeeDeviceDelete   = New("failed to delete device", ErrInternal)
	ErrZigbeeDeviceSearch   = New("failed to search device", ErrInternal)

	ErrUserDeviceGet    = New("failed to get device list", ErrInternal)
	ErrUserDeviceDelete = New("failed to delete user device", ErrInternal)
	ErrUserDeviceAdd    = New("failed to add user device", ErrInternal)

	ErrPolygonAdd    = New("failed to add polygon", ErrInternal)
	ErrPolygonGet    = New("failed to get polygon", ErrInternal)
	ErrPolygonUpdate = New("failed to update polygon", ErrInternal)
)
