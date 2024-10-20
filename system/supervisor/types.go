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

package supervisor

import (
	"context"
	"io/fs"
	"time"

	"github.com/e154/bus"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/web"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
)

// PluginInfo ...
type PluginInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Enabled bool   `json:"enabled"`
	System  bool   `json:"system"`
}

// Supervisor ...
type Supervisor interface {
	Start(context.Context) error
	Shutdown(context.Context) error
	Restart(context.Context) error
	GetPlugin(name string) (interface{}, error)
	EnablePlugin(context.Context, string) error
	DisablePlugin(context.Context, string) error
	PluginList() (list []PluginInfo, total int64, err error)
	SetMetric(common.EntityId, string, map[string]interface{})
	SetState(common.EntityId, EntityStateParams) error
	GetActorById(common.EntityId) (PluginActor, error)
	CallAction(common.EntityId, string, map[string]interface{})
	CallScript(id common.EntityId, fn string, arg ...interface{})
	CallActionV2(CallActionV2, map[string]interface{})
	CallScene(common.EntityId, map[string]interface{})
	AddEntity(*m.Entity) error
	GetEntityById(common.EntityId) (m.EntityShort, error)
	UpdateEntity(*m.Entity) error
	UnloadEntity(common.EntityId)
	EntityIsLoaded(id common.EntityId) bool
	PluginIsLoaded(string) bool
	GetService() Service
	GetPluginReadme(context.Context, string, *string, *string) ([]byte, error)
	PushSystemEvent(strCommand string, params map[string]interface{})
}

// PluginActor ...
type PluginActor interface {
	Spawn()
	Destroy()
	StopWatchers()
	Attributes() m.Attributes
	Settings() m.Attributes
	Metrics() []*m.Metric
	SetState(EntityStateParams) error
	Info() ActorInfo
	GetCurrentState() *events.EventEntityState
	GetOldState() *events.EventEntityState
	SetCurrentState(events.EventEntityState)
	GetEventState() events.EventEntityState
	AddMetric(name string, value map[string]interface{})
	MatchTags(tags []string) bool
	Area() *m.Area
	CallScript(fn string, arg ...interface{})
}

// ActorConstructor ...
type ActorConstructor func(*m.Entity) (PluginActor, error)

// ActorAction ...
type ActorAction struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	ImageUrl     *string                `json:"image_url"`
	Icon         *string                `json:"icon"`
	ScriptEngine *scripts.EngineWatcher `json:"-"`
}

// ToEntityActionShort ...
func ToEntityActionShort(from map[string]ActorAction) (to map[string]m.EntityActionShort) {
	to = make(map[string]m.EntityActionShort)
	for k, v := range from {
		to[k] = m.EntityActionShort{
			Name:        v.Name,
			Description: v.Description,
			ImageUrl:    v.ImageUrl,
			Icon:        v.Icon,
		}
	}
	return
}

// ActorState ...
type ActorState struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageUrl    *string `json:"image_url"`
	Icon        *string `json:"icon"`
}

// ToEntityStateShort ...
func ToEntityStateShort(from map[string]ActorState) (to map[string]m.EntityStateShort) {
	to = make(map[string]m.EntityStateShort)
	for k, v := range from {
		to[k] = m.EntityStateShort{
			Name:        v.Name,
			Description: v.Description,
			ImageUrl:    v.ImageUrl,
			Icon:        v.Icon,
		}
	}
	return
}

// Copy ...
func (a *ActorState) Copy() (state *ActorState) {

	if a == nil {
		return nil
	}

	state = &ActorState{
		Name:        a.Name,
		Description: a.Description,
	}
	if a.ImageUrl != nil {
		state.ImageUrl = common.String(*a.ImageUrl)
	}
	if a.Icon != nil {
		state.Icon = common.String(*a.Icon)
	}
	return
}

const (
	// StateAwait ...
	StateAwait = "await"
	// StateOk ...
	StateOk = "ok"
	// StateError ...
	StateError = "error"
	// StateInProcess ...
	StateInProcess = "in process"
)

// ActorInfo ...
type ActorInfo struct {
	Id                common.EntityId        `json:"id"`
	ParentId          *common.EntityId       `json:"parent_id"`
	PluginName        string                 `json:"plugin_name"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Hidde             bool                   `json:"hidde"`
	UnitOfMeasurement string                 `json:"unit_of_measurement"`
	LastChanged       *time.Time             `json:"last_changed"`
	LastUpdated       *time.Time             `json:"last_updated"`
	DependsOn         []string               `json:"depends_on"`
	State             *ActorState            `json:"state"`
	ImageUrl          *string                `json:"image_url"`
	Icon              *string                `json:"icon"`
	Area              *m.Area                `json:"area"`
	AutoLoad          bool                   `json:"auto_load"`
	RestoreState      bool                   `json:"restoreState"`
	Value             interface{}            `json:"value"`
	States            map[string]ActorState  `json:"states"`
	Actions           map[string]ActorAction `json:"actions"`
}

// PluginType ...
type PluginType string

const (
	// PluginBuiltIn ...
	PluginBuiltIn = PluginType("System")
	// PluginInstallable ...
	PluginInstallable = PluginType("Installable")
)

// Service ...
type Service interface {
	Plugins() map[string]Pluggable
	EventBus() bus.Bus
	Adaptors() *adaptors.Adaptors
	Supervisor() Supervisor
	ScriptService() scripts.ScriptService
	MqttServ() mqtt.MqttServ
	AppConfig() *m.AppConfig
	Scheduler() *scheduler.Scheduler
	Crawler() web.Crawler
}

// Pluggable ...
type Pluggable interface {
	Load(context.Context, Service) error
	Unload(context.Context) error
	Name() string
	Type() PluginType
	Depends() []string
	Version() string
	Options() m.PluginOptions
	EntityIsLoaded(id common.EntityId) bool
	GetActor(id common.EntityId) (pla PluginActor, err error)
	AddOrUpdateActor(*m.Entity) error
	RemoveActor(common.EntityId) error
	Readme(*string, *string) ([]byte, error)
}

// Installable ...
type Installable interface {
	Install() error
	Uninstall() error
}

type CallActionV2 struct {
	EntityId   *common.EntityId `json:"entity_id"`
	ActionName string           `json:"action_name"`
	Tags       []string         `json:"tags"`
	AreaId     *int64           `json:"area_id"`
}

type PluginFileInfo struct {
	ModTime  time.Time
	Name     string
	MimeType string
	Size     int64
	FileMode fs.FileMode
}

type PluginFileInfos []*PluginFileInfo

func (l PluginFileInfos) Len() int      { return len(l) }
func (l PluginFileInfos) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l PluginFileInfos) Less(i, j int) bool {
	return l[i].ModTime.UnixNano() > l[j].ModTime.UnixNano()
}
