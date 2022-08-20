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
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
	"github.com/pkg/errors"
)

// PluginType ...
type PluginType string

const (
	// PluginBuiltIn ...
	PluginBuiltIn = PluginType("System")
	// PluginInstallable ...
	PluginInstallable = PluginType("Installable")
)

var (
	// ErrPluginIsLoaded ...
	ErrPluginIsLoaded = errors.New("plugin is loaded")
	// ErrPluginIsUnloaded ...
	ErrPluginIsUnloaded = errors.New("plugin is unloaded")
	// ErrPluginNotLoaded ...
	ErrPluginNotLoaded = errors.New("plugin not loaded")
)

// Service ...
type Service interface {
	Plugins() map[string]Plugable
	PluginManager() common.PluginManager
	EventBus() bus.Bus
	Adaptors() *adaptors.Adaptors
	EntityManager() entity_manager.EntityManager
	ScriptService() scripts.ScriptService
	MqttServ() mqtt.MqttServ
	AppConfig() *m.AppConfig
	GateClient() *gate_client.GateClient
	Scheduler() *scheduler.Scheduler
}

// Plugable ...
type Plugable interface {
	Load(Service) error
	Unload() error
	Name() string
	Type() PluginType
	Depends() []string
	Version() string
	Options() m.PluginOptions
}

// Installable ...
type Installable interface {
	Install() error
	Uninstall() error
}
