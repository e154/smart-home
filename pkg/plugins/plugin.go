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
	"context"
	"embed"
	"fmt"
	"strings"
	"sync"

	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/logger"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/version"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

var (
	log = logger.MustGetLogger("plugin")
)

// Plugin ...
type Plugin struct {
	Service          Service
	IsStarted        *atomic.Bool
	Actors           *sync.Map
	actorConstructor ActorConstructor
	F                embed.FS
}

// NewPlugin ...
func NewPlugin() *Plugin {
	return &Plugin{
		IsStarted: atomic.NewBool(false),
		Actors:    &sync.Map{},
	}
}

// Load ...
func (p *Plugin) Load(ctx context.Context, service Service, actorConstructor ActorConstructor) error {
	if !p.IsStarted.CompareAndSwap(false, true) {
		return apperr.ErrPluginIsLoaded
	}

	p.Service = service
	p.actorConstructor = actorConstructor

	return nil
}

// Unload ...
func (p *Plugin) Unload(ctx context.Context) error {

	if !p.IsStarted.CompareAndSwap(true, false) {
		return apperr.ErrPluginIsUnloaded
	}

	p.Actors.Range(func(key, value any) bool {
		if pla, ok := value.(PluginActor); ok {
			p.removePluginActor(pla)
		}
		return true
	})

	return nil
}

// Name ...
func (p *Plugin) Name() string {
	panic("implement me")
}

// Depends ...
func (p *Plugin) Depends() []string {
	panic("implement me")
}

// Version ...
func (p *Plugin) Version() string {
	return version.VersionString
}

// Options ...
func (p *Plugin) Options() models.PluginOptions {
	return models.PluginOptions{
		ActorCustomAttrs: false,
		ActorAttrs:       nil,
		ActorSetts:       nil,
	}
}

// LoadSettings ...
func (p *Plugin) LoadSettings(pl Pluggable) (settings models.Attributes, err error) {
	var plugin *models.Plugin
	if plugin, err = p.Service.Adaptors().Plugin.GetByName(context.Background(), pl.Name()); err != nil {
		return
	}
	settings = pl.Options().Setts
	if settings == nil {
		settings = make(models.Attributes)
		return
	}
	_, err = settings.Deserialize(plugin.Settings)
	return
}

func (p *Plugin) AddOrUpdateActor(entity *models.Entity) (err error) {

	if p.actorConstructor == nil {
		return
	}

	var pla PluginActor
	if pla, err = p.actorConstructor(entity); err != nil {
		return
	}

	item, ok := p.Actors.Load(entity.Id)
	if ok && item != nil {
		_ = p.RemoveActor(entity.Id)
	}

	err = p.AddActor(pla, entity)

	return
}

func (p *Plugin) AddActor(pla PluginActor, entity *models.Entity) (err error) {

	if entity == nil {
		return
	}

	p.Actors.Store(entity.Id, pla)
	pla.Spawn()
	log.Infof("entity '%v' loaded", entity.Id)

	p.Service.EventBus().Publish("system/entities/"+entity.Id.String(), events.EventEntityLoaded{
		EntityId:   entity.Id,
		PluginName: entity.PluginName,
	})

	if _, err = p.Service.Adaptors().Entity.GetById(context.Background(), entity.Id); err == nil {
		return
	}

	err = p.Service.Adaptors().Entity.Add(context.Background(), &models.Entity{
		Id:           entity.Id,
		Description:  entity.Description,
		PluginName:   entity.PluginName,
		Icon:         entity.Icon,
		Area:         entity.Area,
		Hidden:       entity.Hidden,
		AutoLoad:     entity.AutoLoad,
		RestoreState: entity.RestoreState,
		ParentId:     entity.ParentId,
		Attributes:   entity.Attributes.Signature(),
		Settings:     entity.Settings,
	})

	return
}

func (p *Plugin) RemoveActor(entityId pkgCommon.EntityId) (err error) {

	item, ok := p.Actors.Load(entityId)
	if !ok || item == nil {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("failed remove \"%s\"", entityId))
		return
	}

	pla := item.(PluginActor)
	p.removePluginActor(pla)
	return
}

func (p *Plugin) removePluginActor(pla PluginActor) {

	info := pla.Info()
	entityId := info.Id

	pla.StopWatchers()
	pla.Destroy()
	p.Actors.Delete(entityId)

	p.Service.EventBus().Publish("system/entities/"+entityId.String(), events.EventEntityUnloaded{
		PluginName: entityId.PluginName(),
		EntityId:   entityId,
	})

	log.Infof("entity '%v' unloaded", entityId)
}

func (p *Plugin) EntityIsLoaded(id pkgCommon.EntityId) bool {
	value, ok := p.Actors.Load(id)
	return ok && value != nil
}

func (p *Plugin) GetActor(id pkgCommon.EntityId) (pla PluginActor, err error) {
	value, ok := p.Actors.Load(id)
	if !ok || value == nil {
		err = errors.Wrap(apperr.ErrEntityNotFound, id.String())
		return
	}
	pla = value.(PluginActor)
	return
}

func (p *Plugin) Readme(note, lang *string) (result []byte, err error) {

	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	n := parser.NewWithExtensions(extensions)

	var fileName = "Readme.md"
	if note != nil {
		fileName = *note
	}
	if !strings.Contains(fileName, ".md") {
		fileName += ".md"
	}
	if lang != nil {
		switch *lang {
		case "ru":
			items := strings.Split(fileName, ".md")
			fileName = strings.Join([]string{items[0], ".", *lang, ".md"}, "")
		default:
		}
	}

	var mds []byte
	if mds, err = p.F.ReadFile(fileName); err != nil {
		return
	}

	doc := n.Parse(mds)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	result = markdown.Render(doc, renderer)

	return
}

func (p *Plugin) Check(msg events.EventCallEntityAction) (result []interface{}, ok bool) {
	if msg.EntityId != nil {
		var value interface{}
		if value, ok = p.Actors.Load(*msg.EntityId); ok {
			result = []interface{}{value}
		}
		return
	}

	p.Actors.Range(func(key, value any) bool {
		pl, _ok := value.(PluginActor)
		if !_ok {
			log.Errorf("error with static cast")
			return false
		}

		var needArea = msg.AreaId != nil
		var needTags = msg.Tags != nil && len(msg.Tags) > 0

		// area
		var areaFound bool
		if needArea {
			areaFound = pl.Area() != nil && pl.Area().Id == *msg.AreaId
			if !needTags && areaFound {
				result = append(result, value)
				return true
			}
		}

		if needArea && !areaFound {
			return true
		}

		// tags
		var tagsFound bool
		if needTags {
			tagsFound = pl.MatchTags(msg.Tags)
			if tagsFound {
				result = append(result, value)
			}
		}

		return true
	})

	ok = len(result) > 0

	return
}
