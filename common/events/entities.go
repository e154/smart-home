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

package events

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
)

// EventStateChanged ...
type EventStateChanged struct {
	StorageSave bool                 `json:"storage_save"`
	PluginName  string               `json:"plugin_name"`
	EntityId    common.EntityId      `json:"entity_id"`
	OldState    bus.EventEntityState `json:"old_state"`
	NewState    bus.EventEntityState `json:"new_state"`
}

// EventLastStateChanged ...
type EventLastStateChanged struct {
	PluginName string               `json:"plugin_name"`
	EntityId   common.EntityId      `json:"entity_id"`
	OldState   bus.EventEntityState `json:"old_state"`
	NewState   bus.EventEntityState `json:"new_state"`
}

// EventGetLastState ...
type EventGetLastState struct {
	EntityId common.EntityId `json:"entity_id"`
}

// EventCallEntityAction ...
type EventCallEntityAction struct {
	PluginName string                 `json:"plugin_name"`
	EntityId   common.EntityId        `json:"entity_id"`
	ActionName string                 `json:"action_name"`
	Args       map[string]interface{} `json:"args"`
}

// EventCallScene ...
type EventCallScene struct {
	PluginName string                 `json:"type"`
	EntityId   common.EntityId        `json:"entity_id"`
	Args       map[string]interface{} `json:"args"`
}

// EventAddedActor ...
//type EventAddedActor struct {
//	PluginName string          `json:"plugin_name"`
//	EntityId   common.EntityId `json:"entity_id"`
//	Attributes m.Attributes    `json:"attributes"`
//	Settings   m.Attributes    `json:"settings"` //???
//}

// EventCreatedEntityModel ...
type EventCreatedEntityModel struct {
	EntityId common.EntityId `json:"entity_id"`
}

// EventUpdatedEntityModel ...
type EventUpdatedEntityModel struct {
	EntityId common.EntityId `json:"entity_id"`
}

// EventUpdatedMetric ...
type EventUpdatedMetric struct {
	EntityId common.EntityId `json:"entity_id"`
}

// CommandUnloadEntity ...
type CommandUnloadEntity struct {
	EntityId common.EntityId `json:"entity_id"`
}

// EventEntityUnloaded ...
type EventEntityUnloaded struct {
	EntityId   common.EntityId `json:"entity_id"`
	PluginName string          `json:"plugin_name"`
}

// CommandLoadEntity ...
type CommandLoadEntity struct {
	EntityId common.EntityId `json:"entity_id"`
}

// EventEntityLoaded ...
type EventEntityLoaded struct {
	EntityId   common.EntityId `json:"entity_id"`
	PluginName string          `json:"plugin_name"`
}

// EventEntitySetState ...
type EventEntitySetState struct {
	EntityId        common.EntityId  `json:"entity_id"`
	NewState        *string          `json:"new_state"`
	AttributeValues m.AttributeValue `json:"attribute_values"`
	SettingsValue   m.AttributeValue `json:"settings_value"`
	StorageSave     bool             `json:"storage_save"`
}
