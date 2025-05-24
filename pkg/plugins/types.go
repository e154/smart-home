// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2024, Filippov Alex
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
	"bufio"
	"context"
	"io/fs"
	"net/http"
	"net/url"
	"time"

	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/models"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/mqtt"
	"github.com/e154/smart-home/pkg/scheduler"
	"github.com/e154/smart-home/pkg/scripts"
	"github.com/e154/smart-home/pkg/web"

	"github.com/e154/bus"
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
	SetMetric(common.EntityId, string, map[string]interface{})
	SetState(common.EntityId, EntityStateParams) error
	GetActorById(common.EntityId) (PluginActor, error)
	CallAction(common.EntityId, string, map[string]interface{})
	CallScript(id common.EntityId, fn string, arg ...interface{})
	CallActionV2(CallActionV2, map[string]interface{})
	CallScene(common.EntityId, map[string]interface{})
	AddEntity(*models.Entity) error
	GetEntityById(common.EntityId) (models.EntityShort, error)
	UpdateEntity(*models.Entity) error
	UnloadEntity(common.EntityId)
	EntityIsLoaded(id common.EntityId) bool
	PluginIsLoaded(string) bool
	GetService() Service
	GetPluginReadme(context.Context, string, *string, *string) ([]byte, error)
	PushSystemEvent(strCommand string, params map[string]interface{})
	UploadPlugin(ctx context.Context, reader *bufio.Reader) (newFile *models.Plugin, err error)
	RemovePlugin(ctx context.Context, pluginName string) error
}

// PluginActor ...
type PluginActor interface {
	Spawn()
	Destroy()
	StopWatchers()
	Attributes() models.Attributes
	Settings() models.Attributes
	Metrics() []*models.Metric
	SetState(EntityStateParams) error
	Info() ActorInfo
	GetCurrentState() *events.EventEntityState
	GetOldState() *events.EventEntityState
	SetCurrentState(events.EventEntityState)
	GetEventState() events.EventEntityState
	AddMetric(name string, value map[string]interface{})
	MatchTags(tags []string) bool
	Area() *models.Area
	CallScript(fn string, arg ...interface{})
}

// ActorConstructor ...
type ActorConstructor func(*models.Entity) (PluginActor, error)

// ActorAction ...
type ActorAction struct {
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	ImageUrl     *string               `json:"image_url"`
	Icon         *string               `json:"icon"`
	ScriptEngine scripts.EngineWatcher `json:"-"`
}

// ToEntityActionShort ...
func ToEntityActionShort(from map[string]ActorAction) (to map[string]models.EntityActionShort) {
	to = make(map[string]models.EntityActionShort)
	for k, v := range from {
		to[k] = models.EntityActionShort{
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
func ToEntityStateShort(from map[string]ActorState) (to map[string]models.EntityStateShort) {
	to = make(map[string]models.EntityStateShort)
	for k, v := range from {
		to[k] = models.EntityStateShort{
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
	Area              *models.Area           `json:"area"`
	AutoLoad          bool                   `json:"auto_load"`
	RestoreState      bool                   `json:"restoreState"`
	Value             interface{}            `json:"value"`
	States            map[string]ActorState  `json:"states"`
	Actions           map[string]ActorAction `json:"actions"`
}

// Service ...
type Service interface {
	Plugins() map[string]Pluggable
	EventBus() bus.Bus
	Adaptors() *adaptors.Adaptors
	Supervisor() Supervisor
	ScriptService() scripts.ScriptService
	MqttServ() mqtt.MqttServ
	AppConfig() *models.AppConfig
	Scheduler() scheduler.Scheduler
	Crawler() web.Crawler
	Authorization() Authorization
	HttpAccessFilter() HttpAccessFilter
}

// Pluggable ...
type Pluggable interface {
	Load(context.Context, Service) error
	Unload(context.Context) error
	Name() string
	Depends() []string
	Version() string
	Options() models.PluginOptions
	EntityIsLoaded(id common.EntityId) bool
	GetActor(id common.EntityId) (pla PluginActor, err error)
	AddOrUpdateActor(*models.Entity) error
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
	IsDir    bool
	FileMode fs.FileMode
}

type PluginFileInfos []*PluginFileInfo

func (l PluginFileInfos) Len() int      { return len(l) }
func (l PluginFileInfos) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l PluginFileInfos) Less(i, j int) bool {
	return l[i].ModTime.UnixNano() > l[j].ModTime.UnixNano()
}

type PluginManifest struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Repository  string                 `json:"repository"`
	Plugin      string                 `json:"plugin"`
	Libs        []string               `json:"libs"`
	Assets      []string               `json:"assets"`
	Settings    map[string]interface{} `json:"settings"`
	Actor       bool                   `json:"actor"`
	Triggers    bool                   `json:"triggers"`
	OS          string                 `json:"os"`
	Arch        string                 `json:"arch"`
}

// EntityStateParams -> supervisor
type EntityStateParams struct {
	NewState        *string          `json:"new_state"`
	AttributeValues m.AttributeValue `json:"attribute_values"`
	SettingsValue   m.AttributeValue `json:"settings_value"`
	StorageSave     bool             `json:"storage_save"`
}

type Authorization interface {
	AuthPlain(login, pass string) (*m.User, error)
	AuthREST(ctx context.Context, accessToken string, requestURI *url.URL, method string) (*m.User, bool, error)
}

type HttpAccessFilter interface {
	Auth(next http.Handler) http.Handler
}
