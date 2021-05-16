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

package plugins

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
)

type PluginType string

const (
	PluginBuiltIn     = PluginType("System")
	PluginInstallable = PluginType("Installable")
)

type Service interface {
	Plugins() map[string]Plugable
	PluginManager() common.PluginManager
	EventBus() event_bus.EventBus
	Adaptors() *adaptors.Adaptors
	EntityManager() entity_manager.EntityManager
	ScriptService() scripts.ScriptService
	MqttServ() mqtt.MqttServ
	AppConfig() *config.AppConfig
	GateClient() *gate_client.GateClient
}

type Plugable interface {
	Load(Service) error
	Unload() error
	Name() string
	Type() PluginType
	Depends() []string
	Version() string
}

type Installable interface {
	Install() error
	Uninstall() error
}
