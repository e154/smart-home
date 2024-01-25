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
	ErrUnknownField    = errors.New("unknown field")
	ErrBadJSONRequest  = errors.New("bad JSON request")
	ErrAccessDenied    = errors.New("access denied")
	ErrAccessForbidden = errors.New("access forbidden")
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
	ErrValidation        = ErrorWithCode("VALIDATION_ERROR", "one or more fields not valid", ErrInvalidRequest)
	ErrInvalidExpiration = ErrorWithCode("INVALID_EXPIRATION", "invalid expiration", ErrInvalidRequest)
	ErrOAuthAccessDenied = ErrorWithCode("ACCESS_DENIED", "client_id or secret not valid", ErrAccessDenied)

	ErrDashboardImport   = ErrorWithCode("DASHBOARD_IMPORT_ERROR", "failed to import dashboard", ErrInternal)
	ErrDashboardAdd      = ErrorWithCode("DASHBOARD_ADD_ERROR", "failed to add dashboard", ErrInternal)
	ErrDashboardGet      = ErrorWithCode("DASHBOARD_GET_ERROR", "failed to get dashboard", ErrInternal)
	ErrDashboardList     = ErrorWithCode("DASHBOARD_LIST_ERROR", "failed to list dashboard", ErrInternal)
	ErrDashboardNotFound = ErrorWithCode("DASHBOARD_NOT_FOUND_ERROR", "dashboard is not found", ErrNotFound)
	ErrDashboardUpdate   = ErrorWithCode("DASHBOARD_UPDATE_ERROR", "failed to update dashboard", ErrInternal)
	ErrDashboardDelete   = ErrorWithCode("DASHBOARD_DELETE_ERROR", "failed to delete dashboard", ErrInternal)
	ErrDashboardSearch   = ErrorWithCode("DASHBOARD_SEARCH_ERROR", "failed to search dashboard", ErrInternal)

	ErrDashboardCardAdd      = ErrorWithCode("DASHBOARD_CARD_ADD_ERROR", "failed to add dashboard card", ErrInternal)
	ErrDashboardCardGet      = ErrorWithCode("DASHBOARD_CARD_GET_ERROR", "failed to get dashboard card", ErrInternal)
	ErrDashboardCardList     = ErrorWithCode("DASHBOARD_CARD_LIST_ERROR", "failed to list dashboard card", ErrInternal)
	ErrDashboardCardNotFound = ErrorWithCode("DASHBOARD_CARD_NOT_FOUND_ERROR", "dashboard card is not found", ErrNotFound)
	ErrDashboardCardUpdate   = ErrorWithCode("DASHBOARD_CARD_UPDATE_ERROR", "failed to update dashboard card", ErrInternal)
	ErrDashboardCardDelete   = ErrorWithCode("DASHBOARD_CARD_DELETE_ERROR", "failed to delete dashboard card", ErrInternal)

	ErrDashboardCardItemAdd      = ErrorWithCode("DASHBOARD_CARD_ITEM_ADD_ERROR", "failed to add dashboard card item", ErrInternal)
	ErrDashboardCardItemGet      = ErrorWithCode("DASHBOARD_CARD_ITEM_GET_ERROR", "failed to get dashboard card item", ErrInternal)
	ErrDashboardCardItemList     = ErrorWithCode("DASHBOARD_CARD_ITEM_LIST_ERROR", "failed to list dashboard card item", ErrInternal)
	ErrDashboardCardItemNotFound = ErrorWithCode("DASHBOARD_CARD_ITEM_NOT_FOUND_ERROR", "dashboard card item is not found", ErrNotFound)
	ErrDashboardCardItemUpdate   = ErrorWithCode("DASHBOARD_CARD_ITEM_UPDATE_ERROR", "failed to update dashboard card item", ErrInternal)
	ErrDashboardCardItemDelete   = ErrorWithCode("DASHBOARD_CARD_ITEM_DELETE_ERROR", "failed to delete dashboard card item", ErrInternal)

	ErrDashboardTabAdd      = ErrorWithCode("DASHBOARD_TAB_ADD_ERROR", "failed to add dashboard tab", ErrInternal)
	ErrDashboardTabGet      = ErrorWithCode("DASHBOARD_TAB_GET_ERROR", "failed to get dashboard tab", ErrInternal)
	ErrDashboardTabList     = ErrorWithCode("DASHBOARD_TAB_LIST_ERROR", "failed to list dashboard tab", ErrInternal)
	ErrDashboardTabNotFound = ErrorWithCode("DASHBOARD_TAB_NOT_FOUND_ERROR", "dashboard tab is not found", ErrNotFound)
	ErrDashboardTabUpdate   = ErrorWithCode("DASHBOARD_TAB_UPDATE_ERROR", "failed to update dashboard tab", ErrInternal)
	ErrDashboardTabDelete   = ErrorWithCode("DASHBOARD_TAB_DELETE_ERROR", "failed to delete dashboard tab", ErrInternal)

	ErrActionAdd      = ErrorWithCode("ACTION_ADD_ERROR", "failed to add action", ErrInternal)
	ErrActionGet      = ErrorWithCode("ACTION_GET_ERROR", "failed to get action", ErrInternal)
	ErrActionUpdate   = ErrorWithCode("ACTION_UPDATE_ERROR", "failed to update action", ErrInternal)
	ErrActionList     = ErrorWithCode("ACTION_LIST_ERROR", "failed to list action", ErrInternal)
	ErrActionNotFound = ErrorWithCode("ACTION_NOT_FOUND_ERROR", "action is not found", ErrNotFound)
	ErrActionDelete   = ErrorWithCode("ACTION_DELETE_ERROR", "failed to delete action", ErrInternal)
	ErrActionSearch   = ErrorWithCode("ACTION_SEARCH_ERROR", "failed to search action", ErrInternal)

	ErrEntityAdd           = ErrorWithCode("ENTITY_ADD_ERROR", "failed to add entity", ErrInternal)
	ErrEntityGet           = ErrorWithCode("ENTITY_GET_ERROR", "failed to get entity", ErrInternal)
	ErrEntityList          = ErrorWithCode("ENTITY_LIST_ERROR", "failed to list entity", ErrInternal)
	ErrEntityNotFound      = ErrorWithCode("ENTITY_NOT_FOUND_ERROR", "entity is not found", ErrNotFound)
	ErrEntityUpdate        = ErrorWithCode("ENTITY_UPDATE_ERROR", "failed to update entity", ErrInternal)
	ErrEntityDelete        = ErrorWithCode("ENTITY_DELETE_ERROR", "failed to delete entity", ErrInternal)
	ErrEntitySerch         = ErrorWithCode("ENTITY_SERCH_ERROR", "failed to search entity", ErrInternal)
	ErrEntityAppendMetric  = ErrorWithCode("ENTITY_APPEND_METRIC_ERROR", "entity append metric", ErrInternal)
	ErrEntityDeleteMetric  = ErrorWithCode("ENTITY_DELETE_METRIC_ERROR", "delete metric failed", ErrInternal)
	ErrEntityReplaceMetric = ErrorWithCode("ENTITY_REPLACE_METRIC_ERROR", "replace metric failed", ErrInternal)
	ErrEntityAppendScript  = ErrorWithCode("ENTITY_APPEND_SCRIPT_ERROR", "append script failed", ErrInternal)
	ErrEntityDeleteScript  = ErrorWithCode("ENTITY_DELETE_SCRIPT_ERROR", "delete script failed", ErrInternal)
	ErrEntityReplaceScript = ErrorWithCode("ENTITY_REPLACE_SCRIPT_ERROR", "replace script failed", ErrInternal)

	ErrAlexaIntentAdd      = ErrorWithCode("ALEXA_INTENT_ADD_ERROR", "failed to add intent", ErrInternal)
	ErrAlexaIntentUpdate   = ErrorWithCode("ALEXA_INTENT_UPDATE_ERROR", "failed to update intent", ErrInternal)
	ErrAlexaIntentGet      = ErrorWithCode("ALEXA_INTENT_GET_ERROR", "failed to get intent", ErrInternal)
	ErrAlexaIntentDelete   = ErrorWithCode("ALEXA_INTENT_DELETE_ERROR", "failed to delete intent", ErrInternal)
	ErrAlexaIntentNotFound = ErrorWithCode("ALEXA_INTENT_NOT_FOUND_ERROR", "intent is not found", ErrNotFound)

	ErrAlexaSkillAdd      = ErrorWithCode("ALEXA_SKILL_ADD_ERROR", "failed to add skill", ErrInternal)
	ErrAlexaSkillGet      = ErrorWithCode("ALEXA_SKILL_GET_ERROR", "failed to get skill", ErrInternal)
	ErrAlexaSkillUpdate   = ErrorWithCode("ALEXA_SKILL_UPDATE_ERROR", "failed to update skill", ErrInternal)
	ErrAlexaSkillList     = ErrorWithCode("ALEXA_SKILL_LIST_ERROR", "failed to list skill", ErrInternal)
	ErrAlexaSkillNotFound = ErrorWithCode("ALEXA_SKILL_NOT_FOUND_ERROR", "skill is not found", ErrNotFound)
	ErrAlexaSkillDelete   = ErrorWithCode("ALEXA_SKILL_DELETE_ERROR", "failed to delete skill", ErrInternal)

	ErrAreaAdd      = ErrorWithCode("AREA_ADD_ERROR", "failed to add area", ErrInternal)
	ErrAreaGet      = ErrorWithCode("AREA_GET_ERROR", "failed to get area", ErrInternal)
	ErrAreaUpdate   = ErrorWithCode("AREA_UPDATE_ERROR", "failed to update area", ErrInternal)
	ErrAreaList     = ErrorWithCode("AREA_LIST_ERROR", "failed to list area", ErrInternal)
	ErrAreaNotFound = ErrorWithCode("AREA_NOT_FOUND_ERROR", "area is not found", ErrNotFound)
	ErrAreaDelete   = ErrorWithCode("AREA_DELETE_ERROR", "failed to delete area", ErrInternal)
	ErrAreaClean    = ErrorWithCode("AREA_CLEAN_ERROR", "failed to clean area", ErrInternal)

	ErrConditionAdd      = ErrorWithCode("CONDITION_ADD_ERROR", "failed to add condition", ErrInternal)
	ErrConditionGet      = ErrorWithCode("CONDITION_GET_ERROR", "failed to get condition", ErrInternal)
	ErrConditionUpdate   = ErrorWithCode("CONDITION_UPDATE_ERROR", "failed to update condition", ErrInternal)
	ErrConditionList     = ErrorWithCode("CONDITION_LIST_ERROR", "failed to list condition", ErrInternal)
	ErrConditionNotFound = ErrorWithCode("CONDITION_NOT_FOUND_ERROR", "condition is not found", ErrNotFound)
	ErrConditionDelete   = ErrorWithCode("CONDITION_DELETE_ERROR", "failed to delete condition", ErrInternal)
	ErrConditionSearch   = ErrorWithCode("CONDITION_SEARCH_ERROR", "failed to search condition", ErrInternal)

	ErrEntityActionAdd      = ErrorWithCode("ENTITY_ACTION_ADD_ERROR", "failed to add action", ErrInternal)
	ErrEntityActionGet      = ErrorWithCode("ENTITY_ACTION_GET_ERROR", "failed to get action", ErrInternal)
	ErrEntityActionUpdate   = ErrorWithCode("ENTITY_ACTION_UPDATE_ERROR", "failed to update action", ErrInternal)
	ErrEntityActionList     = ErrorWithCode("ENTITY_ACTION_LIST_ERROR", "failed to list action", ErrInternal)
	ErrEntityActionNotFound = ErrorWithCode("ENTITY_ACTION_NOT_FOUND_ERROR", "action is not found", ErrNotFound)
	ErrEntityActionDelete   = ErrorWithCode("ENTITY_ACTION_DELETE_ERROR", "failed to delete action", ErrInternal)

	ErrEntityStateAdd      = ErrorWithCode("ENTITY_STATE_ADD_ERROR", "failed to add state", ErrInternal)
	ErrEntityStateGet      = ErrorWithCode("ENTITY_STATE_GET_ERROR", "failed to get state", ErrInternal)
	ErrEntityStateUpdate   = ErrorWithCode("ENTITY_STATE_UPDATE_ERROR", "failed to update state", ErrInternal)
	ErrEntityStateList     = ErrorWithCode("ENTITY_STATE_LIST_ERROR", "failed to list state", ErrInternal)
	ErrEntityStateNotFound = ErrorWithCode("ENTITY_STATE_NOT_FOUND_ERROR", "state is not found", ErrNotFound)
	ErrEntityStateDelete   = ErrorWithCode("ENTITY_STATE_DELETE_ERROR", "failed to delete state", ErrInternal)

	ErrEntityStorageAdd    = ErrorWithCode("ENTITY_STORAGE_ADD_ERROR", "failed to add storage", ErrInternal)
	ErrEntityStorageGet    = ErrorWithCode("ENTITY_STORAGE_GET_ERROR", "failed to get storage", ErrInternal)
	ErrEntityStorageList   = ErrorWithCode("ENTITY_STORAGE_LIST_ERROR", "failed to list storage", ErrInternal)
	ErrEntityStorageDelete = ErrorWithCode("ENTITY_STORAGE_DELETE_ERROR", "failed to delete storage", ErrInternal)

	ErrImageAdd      = ErrorWithCode("IMAGE_ADD_ERROR", "failed to add image", ErrInternal)
	ErrImageGet      = ErrorWithCode("IMAGE_GET_ERROR", "failed to get image", ErrInternal)
	ErrImageUpdate   = ErrorWithCode("IMAGE_UPDATE_ERROR", "failed to update image", ErrInternal)
	ErrImageList     = ErrorWithCode("IMAGE_LIST_ERROR", "failed to list image", ErrInternal)
	ErrImageNotFound = ErrorWithCode("IMAGE_NOT_FOUND_ERROR", "image is not found", ErrNotFound)
	ErrImageDelete   = ErrorWithCode("IMAGE_DELETE_ERROR", "failed to delete image", ErrInternal)

	ErrLogAdd      = ErrorWithCode("LOG_ADD_ERROR", "failed to add log", ErrInternal)
	ErrLogGet      = ErrorWithCode("LOG_GET_ERROR", "failed to get log", ErrInternal)
	ErrLogList     = ErrorWithCode("LOG_LIST_ERROR", "failed to list log", ErrInternal)
	ErrLogNotFound = ErrorWithCode("LOG_NOT_FOUND_ERROR", "log is not found", ErrNotFound)
	ErrLogDelete   = ErrorWithCode("LOG_DELETE_ERROR", "failed to delete log", ErrNotFound)

	ErrMessageAdd              = ErrorWithCode("MESSAGE_ADD_ERROR", "failed to add message", ErrInternal)
	ErrMessageDeliveryAdd      = ErrorWithCode("MESSAGE_DELIVERY_ADD_ERROR", "failed to add message delivery", ErrInternal)
	ErrMessageDeliveryList     = ErrorWithCode("MESSAGE_DELIVERY_LIST_ERROR", "failed to list message delivery", ErrInternal)
	ErrMessageDeliveryUpdate   = ErrorWithCode("MESSAGE_DELIVERY_UPDATE_ERROR", "failed to update message delivery", ErrInternal)
	ErrMessageDeliveryDelete   = ErrorWithCode("MESSAGE_DELIVERY_DELETE_ERROR", "failed to delete message delivery", ErrInternal)
	ErrMessageDeliveryGet      = ErrorWithCode("MESSAGE_DELIVERY_GET_ERROR", "failed to get message delivery", ErrInternal)
	ErrMessageDeliveryNotFound = ErrorWithCode("MESSAGE_DELIVERY_NOT_FOUND_ERROR", "message delivery is not found", ErrNotFound)

	ErrMetricAdd      = ErrorWithCode("METRIC_ADD_ERROR", "failed to add metric", ErrInternal)
	ErrMetricGet      = ErrorWithCode("METRIC_GET_ERROR", "failed to get metric", ErrInternal)
	ErrMetricUpdate   = ErrorWithCode("METRIC_UPDATE_ERROR", "failed to update metric", ErrInternal)
	ErrMetricList     = ErrorWithCode("METRIC_LIST_ERROR", "failed to list metric", ErrInternal)
	ErrMetricNotFound = ErrorWithCode("METRIC_NOT_FOUND_ERROR", "metric is not found", ErrNotFound)
	ErrMetricDelete   = ErrorWithCode("METRIC_DELETE_ERROR", "failed to delete metric", ErrInternal)
	ErrMetricSearch   = ErrorWithCode("METRIC_SEARCH_ERROR", "failed to search metric", ErrInternal)

	ErrMetricBucketAdd    = ErrorWithCode("METRIC_BUCKET_ADD_ERROR", "failed to add metric backet", ErrInternal)
	ErrMetricBucketGet    = ErrorWithCode("METRIC_BUCKET_GET_ERROR", "failed to get metric backet", ErrInternal)
	ErrMetricBucketDelete = ErrorWithCode("METRIC_BUCKET_DELETE_ERROR", "failed to delete metric backet", ErrInternal)

	ErrPermissionAdd    = ErrorWithCode("PERMISSION_ADD_ERROR", "failed to add permission", ErrInternal)
	ErrPermissionGet    = ErrorWithCode("PERMISSION_GET_ERROR", "failed to get permission", ErrInternal)
	ErrPermissionDelete = ErrorWithCode("PERMISSION_DELETE_ERROR", "failed to delete permission", ErrInternal)

	ErrPluginAdd        = ErrorWithCode("PLUGIN_ADD_ERROR", "failed to add plugin", ErrInternal)
	ErrPluginGet        = ErrorWithCode("PLUGIN_GET_ERROR", "failed to get plugin", ErrInternal)
	ErrPluginUpdate     = ErrorWithCode("PLUGIN_UPDATE_ERROR", "failed to update plugin", ErrInternal)
	ErrPluginList       = ErrorWithCode("PLUGIN_LIST_ERROR", "failed to list plugin", ErrInternal)
	ErrPluginNotFound   = ErrorWithCode("PLUGIN_NOT_FOUND_ERROR", "plugin is not found", ErrNotFound)
	ErrPluginDelete     = ErrorWithCode("PLUGIN_DELETE_ERROR", "failed to delete plugin", ErrInternal)
	ErrPluginSearch     = ErrorWithCode("PLUGIN_SEARCH_ERROR", "failed to search plugin", ErrInternal)
	ErrPluginIsLoaded   = ErrorWithCode("PLUGIN_IS_LOADED", "plugin is loaded", ErrInvalidRequest)
	ErrPluginIsUnloaded = ErrorWithCode("PLUGIN_IS_UNLOADED", "plugin is unloaded", ErrInvalidRequest)
	ErrPluginNotLoaded  = ErrorWithCode("PLUGIN_NOT_LOADED", "plugin not loaded", ErrInvalidRequest)

	ErrRoleAdd             = ErrorWithCode("ROLE_ADD_ERROR", "failed to add role", ErrInternal)
	ErrRoleGet             = ErrorWithCode("ROLE_GET_ERROR", "failed to get role", ErrInternal)
	ErrRoleUpdate          = ErrorWithCode("ROLE_UPDATE_ERROR", "failed to update role", ErrInternal)
	ErrRoleUpdateForbidden = ErrorWithCode("ROLE_UPDATE_ERROR", "failed to update role", ErrAccessForbidden)
	ErrRoleList            = ErrorWithCode("ROLE_LIST_ERROR", "failed to list role", ErrInternal)
	ErrRoleNotFound        = ErrorWithCode("ROLE_NOT_FOUND_ERROR", "role is not found", ErrNotFound)
	ErrRoleDelete          = ErrorWithCode("ROLE_DELETE_ERROR", "failed to delete role", ErrInternal)
	ErrRoleDeleteForbidden = ErrorWithCode("ROLE_DELETE_ERROR", "failed to delete role", ErrAccessForbidden)
	ErrRoleSearch          = ErrorWithCode("ROLE_SEARCH_ERROR", "failed to search role", ErrInternal)

	ErrRunStoryAdd    = ErrorWithCode("RUN_STORY_ADD_ERROR", "failed to add run story", ErrInternal)
	ErrRunStoryUpdate = ErrorWithCode("RUN_STORY_UPDATE_ERROR", "failed to update run story", ErrInternal)
	ErrRunStoryList   = ErrorWithCode("RUN_STORY_LIST_ERROR", "failed to list run story", ErrInternal)

	ErrScriptAdd      = ErrorWithCode("SCRIPT_ADD_ERROR", "failed to add script", ErrInternal)
	ErrScriptGet      = ErrorWithCode("SCRIPT_GET_ERROR", "failed to get script", ErrInternal)
	ErrScriptUpdate   = ErrorWithCode("SCRIPT_UPDATE_ERROR", "failed to update script", ErrInternal)
	ErrScriptList     = ErrorWithCode("SCRIPT_LIST_ERROR", "failed to list script", ErrInternal)
	ErrScriptNotFound = ErrorWithCode("SCRIPT_NOT_FOUND_ERROR", "script is not found", ErrNotFound)
	ErrScriptDelete   = ErrorWithCode("SCRIPT_DELETE_ERROR", "failed to delete script", ErrInternal)
	ErrScriptSearch   = ErrorWithCode("SCRIPT_SEARCH_ERROR", "failed to search script", ErrInternal)
	ErrScriptStat     = ErrorWithCode("SCRIPT_STAT_ERROR", "failed to get script statistic", ErrInternal)

	ErrTaskAdd             = ErrorWithCode("TASK_ADD_ERROR", "failed to add Task", ErrInternal)
	ErrTaskGet             = ErrorWithCode("TASK_GET_ERROR", "failed to get Task", ErrInternal)
	ErrTaskUpdate          = ErrorWithCode("TASK_UPDATE_ERROR", "failed to update Task", ErrInternal)
	ErrTaskList            = ErrorWithCode("TASK_LIST_ERROR", "failed to list Task", ErrInternal)
	ErrTaskNotFound        = ErrorWithCode("TASK_NOT_FOUND_ERROR", "Task is not found", ErrNotFound)
	ErrTaskDelete          = ErrorWithCode("TASK_DELETE_ERROR", "failed to delete Task", ErrInternal)
	ErrTaskSearch          = ErrorWithCode("TASK_SEARCH_ERROR", "failed to search Task", ErrInternal)
	ErrTaskAppendTrigger   = ErrorWithCode("TASK_APPEND_TRIGGER_ERROR", "task append trigger failed", ErrInternal)
	ErrTaskDeleteTrigger   = ErrorWithCode("TASK_DELETE_TRIGGER_ERROR", "task delete trigger failed", ErrInternal)
	ErrTaskAppendCondition = ErrorWithCode("TASK_APPEND_CONDITION_ERROR", "task append condition failed", ErrInternal)
	ErrTaskDeleteCondition = ErrorWithCode("TASK_DELETE_CONDITION_ERROR", "task delete condition failed", ErrInternal)
	ErrTaskAppendAction    = ErrorWithCode("TASK_APPEND_ACTION_ERROR", "task append action failed", ErrInternal)
	ErrTaskDeleteAction    = ErrorWithCode("TASK_DELETE_ACTION_ERROR", "task delete action failed", ErrInternal)

	ErrChatAdd    = ErrorWithCode("CHAT_ADD_ERROR", "failed to add chat", ErrInternal)
	ErrChatList   = ErrorWithCode("CHAT_LIST_ERROR", "failed to list chat", ErrInternal)
	ErrChatDelete = ErrorWithCode("CHAT_DELETE_ERROR", "failed to delete chat", ErrInternal)

	ErrTemplateAdd      = ErrorWithCode("TEMPLATE_ADD_ERROR", "failed to add template", ErrInternal)
	ErrTemplateGet      = ErrorWithCode("TEMPLATE_GET_ERROR", "failed to get template", ErrInternal)
	ErrTemplateUpdate   = ErrorWithCode("TEMPLATE_UPDATE_ERROR", "failed to update template", ErrInternal)
	ErrTemplateList     = ErrorWithCode("TEMPLATE_LIST_ERROR", "failed to list template", ErrInternal)
	ErrTemplateNotFound = ErrorWithCode("TEMPLATE_NOT_FOUND_ERROR", "template is not found", ErrNotFound)
	ErrTemplateDelete   = ErrorWithCode("TEMPLATE_DELETE_ERROR", "failed to delete template", ErrInternal)
	ErrTemplateSearch   = ErrorWithCode("TEMPLATE_SEARCH_ERROR", "failed to search template", ErrInternal)

	ErrTriggerAdd          = ErrorWithCode("TRIGGER_ADD_ERROR", "failed to add trigger", ErrInternal)
	ErrTriggerGet          = ErrorWithCode("TRIGGER_GET_ERROR", "failed to get trigger", ErrInternal)
	ErrTriggerUpdate       = ErrorWithCode("TRIGGER_UPDATE_ERROR", "failed to update trigger", ErrInternal)
	ErrTriggerList         = ErrorWithCode("TRIGGER_LIST_ERROR", "failed to list trigger", ErrInternal)
	ErrTriggerNotFound     = ErrorWithCode("TRIGGER_NOT_FOUND_ERROR", "trigger is not found", ErrNotFound)
	ErrTriggerDelete       = ErrorWithCode("TRIGGER_DELETE_ERROR", "failed to delete trigger", ErrInternal)
	ErrTriggerSearch       = ErrorWithCode("TRIGGER_SEARCH_ERROR", "failed to search trigger", ErrInternal)
	ErrTriggerDeleteEntity = ErrorWithCode("TRIGGER_DELETE_ENTITY_ERROR", "trigger delete entity failed", ErrInternal)

	ErrUserAdd             = ErrorWithCode("USER_ADD_ERROR", "failed to add user", ErrInternal)
	ErrUserMetaAdd         = ErrorWithCode("USER_META_ADD_ERROR", "failed to add user meta", ErrInternal)
	ErrUserGet             = ErrorWithCode("USER_GET_ERROR", "failed to get user", ErrInternal)
	ErrUserUpdate          = ErrorWithCode("USER_UPDATE_ERROR", "failed to update user", ErrInternal)
	ErrUserUpdateForbidden = ErrorWithCode("USER_UPDATE_ERROR", "failed to update user", ErrAccessForbidden)
	ErrUserList            = ErrorWithCode("USER_LIST_ERROR", "failed to list user", ErrInternal)
	ErrUserNotFound        = ErrorWithCode("USER_NOT_FOUND_ERROR", "user is not found", ErrNotFound)
	ErrUserDelete          = ErrorWithCode("USER_DELETE_ERROR", "failed to delete user", ErrInternal)
	ErrUserDeleteForbidden = ErrorWithCode("USER_DELETE_ERROR", "failed to delete user", ErrAccessForbidden)

	ErrVariableAdd             = ErrorWithCode("VARIABLE_ADD_ERROR", "failed to add variable", ErrInternal)
	ErrVariableGet             = ErrorWithCode("VARIABLE_GET_ERROR", "failed to get variable", ErrInternal)
	ErrVariableUpdate          = ErrorWithCode("VARIABLE_UPDATE_ERROR", "failed to update variable", ErrInternal)
	ErrVariableList            = ErrorWithCode("VARIABLE_LIST_ERROR", "failed to list variable", ErrInternal)
	ErrVariableNotFound        = ErrorWithCode("VARIABLE_NOT_FOUND_ERROR", "variable is not found", ErrNotFound)
	ErrVariableDelete          = ErrorWithCode("VARIABLE_DELETE_ERROR", "failed to delete variable", ErrInternal)
	ErrVariableUpdateForbidden = ErrorWithCode("VARIABLE_UPDATE_ERROR", "unable to update system variable", ErrAccessForbidden)

	ErrZigbee2mqttAdd      = ErrorWithCode("ZIGBEE2MQTT_ADD_ERROR", "failed to add zigbee2mqtt", ErrInternal)
	ErrZigbee2mqttGet      = ErrorWithCode("ZIGBEE2MQTT_GET_ERROR", "failed to get zigbee2mqtt", ErrInternal)
	ErrZigbee2mqttUpdate   = ErrorWithCode("ZIGBEE2MQTT_UPDATE_ERROR", "failed to update zigbee2mqtt", ErrInternal)
	ErrZigbee2mqttList     = ErrorWithCode("ZIGBEE2MQTT_LIST_ERROR", "failed to list zigbee2mqtt", ErrInternal)
	ErrZigbee2mqttNotFound = ErrorWithCode("ZIGBEE2MQTT_NOT_FOUND_ERROR", "zigbee2mqtt is not found", ErrNotFound)
	ErrZigbee2mqttDelete   = ErrorWithCode("ZIGBEE2MQTT_DELETE_ERROR", "failed to delete zigbee2mqtt", ErrInternal)

	ErrZigbeeDeviceAdd      = ErrorWithCode("ZIGBEE_DEVICE_ADD_ERROR", "failed to add device", ErrInternal)
	ErrZigbeeDeviceGet      = ErrorWithCode("ZIGBEE_DEVICE_GET_ERROR", "failed to get device", ErrInternal)
	ErrZigbeeDeviceUpdate   = ErrorWithCode("ZIGBEE_DEVICE_UPDATE_ERROR", "failed to update device", ErrInternal)
	ErrZigbeeDeviceList     = ErrorWithCode("ZIGBEE_DEVICE_LIST_ERROR", "failed to list device", ErrInternal)
	ErrZigbeeDeviceNotFound = ErrorWithCode("ZIGBEE_DEVICE_NOT_FOUND_ERROR", "device is not found", ErrNotFound)
	ErrZigbeeDeviceDelete   = ErrorWithCode("ZIGBEE_DEVICE_DELETE_ERROR", "failed to delete device", ErrInternal)
	ErrZigbeeDeviceSearch   = ErrorWithCode("ZIGBEE_DEVICE_SEARCH_ERROR", "failed to search device", ErrInternal)

	ErrUserDeviceGet    = ErrorWithCode("USER_DEVICE_GET_ERROR", "failed to get device list", ErrInternal)
	ErrUserDeviceDelete = ErrorWithCode("USER_DEVICE_DELETE_ERROR", "failed to delete user device", ErrInternal)
	ErrUserDeviceAdd    = ErrorWithCode("USER_DEVICE_ADD_ERROR", "failed to add user device", ErrInternal)

	ErrBackupNotFound           = ErrorWithCode("BACKUP_NOT_FOUND_ERROR", "backup not found", ErrNotFound)
	ErrBackupNameNotUnique      = ErrorWithCode("BACKUP_NAME_NOT_UNIQUE_ERROR", "backup name not unique", ErrInvalidRequest)
	ErrBackupRestoreForbidden   = ErrorWithCode("BACKUP_RESTORE_ERROR", "failed to restore backup", ErrAccessForbidden)
	ErrBackupApplyForbidden     = ErrorWithCode("BACKUP_APPLY_ERROR", "failed to apply backup", ErrAccessForbidden)
	ErrBackupRollbackForbidden  = ErrorWithCode("BACKUP_ROLLBACK_ERROR", "failed to rollback backup", ErrAccessForbidden)
	ErrBackupCreateNewForbidden = ErrorWithCode("BACKUP_CREATE_ERROR", "failed to create new backup", ErrAccessForbidden)
	ErrBackupUploadForbidden    = ErrorWithCode("BACKUP_UPLOAD_ERROR", "failed to upload backup", ErrAccessForbidden)

	ErrScriptVersionAdd    = ErrorWithCode("SCRIPT_VERSION_ADD_ERROR", "failed to add script version", ErrInternal)
	ErrScriptVersionList   = ErrorWithCode("SCRIPT_VERSION_LIST_ERROR", "failed to list script version", ErrInternal)
	ErrScriptVersionDelete = ErrorWithCode("SCRIPT_VERSION_DELETE_ERROR", "failed to delete script version", ErrInternal)
)
