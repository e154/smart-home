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

package plugin

import (
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/e154/smart-home/common/logger"

	"github.com/e154/smart-home/cmd/cli/commands/generate"
	"github.com/e154/smart-home/common"
	"github.com/spf13/cobra"
)

var (
	log = logger.MustGetLogger("plugin")
)

var actorTpl = `//CODE GENERATED AUTOMATICALLY

package {{.Package}}

import (
	"fmt"
	"sync"

	"{{.Dir}}/common"
	m "{{.Dir}}/models"
	"{{.Dir}}/system/bus"
	"{{.Dir}}/system/entity_manager"
)

type Actor struct {
	entity_manager.BaseActor
	eventBus bus.Bus
}

func NewActor(entity *m.Entity,
	entityManager entity_manager.EntityManager,
	eventBus bus.Bus) *Actor {

	name := entity.Id.Name()

	actor := &Actor{
		BaseActor: entity_manager.BaseActor{
			Id:          common.EntityId(fmt.Sprintf("%s.%s", Entity{{.PluginName}}, name)),
			Name:        name,
			Description: "{{.PluginName}} plugin",
			EntityType:  Entity{{.PluginName}},
			AttrMu:      &sync.RWMutex{},
			Attrs:       entity.Attributes,
			Setts:       entity.Settings,
			Manager:     entityManager,
			States:      NewStates(),
			Actions:     NewActions(),
		},
		eventBus: eventBus,
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	return actor
}

func (e *Actor) Spawn() entity_manager.PluginActor {
	e.Update()
	return e
}

func (e *Actor) Update() {

}

`

var pluginTpl = `//CODE GENERATED AUTOMATICALLY

package {{.Package}}

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"{{.Dir}}/common"
	"{{.Dir}}/common/apperr"
	"{{.Dir}}/common/logger"
	m "{{.Dir}}/models"
	"{{.Dir}}/system/entity_manager"
	"{{.Dir}}/system/plugins"
)

var (
	log = logger.MustGetLogger("plugins.{{.Package}}")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	actorsLock *sync.Mutex
	actors     map[common.EntityId]*Actor
	quit       chan struct{}
	pause      time.Duration
}

func New() plugins.Plugable {
	return &plugin{
		Plugin:     plugins.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[common.EntityId]*Actor),
		pause:      240,
	}
}

func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.quit = make(chan struct{})

	go func() {
		ticker := time.NewTicker(time.Second * p.pause)

		defer func() {
			ticker.Stop()
			close(p.quit)
		}()

		for {
			select {
			case <-p.quit:
				return
			case <-ticker.C:
				p.updateForAll()
			}
		}
	}()

	return nil
}

func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	p.quit <- struct{}{}
	return nil
}

func (p *plugin) Name() string {
	return Name
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entity.Id]; ok {
		p.actors[entity.Id].Update()
		return
	}

	p.actors[entity.Id] = NewActor(entity, p.EntityManager, p.EventBus)
	p.EntityManager.Spawn(p.actors[entity.Id].Spawn)

	return
}

func (p *plugin) RemoveActor(entityId common.EntityId) error {
	return p.removeEntity(entityId)
}

func (p *plugin) removeEntity(entityId common.EntityId) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entityId]; !ok {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("failed remove \"%s\"", entityId.Name()))
		return
	}

	delete(p.actors, entityId)

	return
}

func (p *plugin) updateForAll() {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	for _, actor := range p.actors {
		actor.Update()
	}
}

func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginBuiltIn
}

func (p *plugin) Depends() []string {
	return nil
}

func (p *plugin) Version() string {
	return "0.0.1"
}

func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Actors:       true,
		ActorAttrs:   NewAttr(),
		ActorSetts:   NewSettings(),
		ActorActions: entity_manager.ToEntityActionShort(NewActions()),
		ActorStates:  entity_manager.ToEntityStateShort(NewStates()),
	}
}

`

var typesTpl = `//CODE GENERATED AUTOMATICALLY

package {{.Package}}

import (
	"{{.Dir}}/common"
	m "{{.Dir}}/models"
	"{{.Dir}}/system/entity_manager"
)

const (
	Name        = "{{.PluginName}}"
	Entity{{.PluginName}} = string("{{.PluginName}}")
)

const (
	SettingParam1 = "param1"
	SettingParam2 = "param2"

	AttrPhase = "phase"

	StateEnabled  = "enabled"
	StateDisabled = "disabled"

	ActionEnabled = "enable"
	ActionDisable = "disable"
)

// store entity status in this struct
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrPhase: {
			Name: AttrPhase,
			Type: common.AttributeString,
		},
	}
}

// entity settings
func NewSettings() m.Attributes {
	return m.Attributes{
		SettingParam1: {
			Name: SettingParam1,
			Type: common.AttributeString,
		},
		SettingParam2: {
			Name: SettingParam2,
			Type: common.AttributeString,
		},
	}
}

// state list entity
func NewStates() (states map[string]entity_manager.ActorState) {

	states = map[string]entity_manager.ActorState{
		StateEnabled: {
			Name:        StateEnabled,
			Description: "Enabled",
		},
		StateDisabled: {
			Name:        StateDisabled,
			Description: "disabled",
		},
	}

	return
}

// entity action list
func NewActions() map[string]entity_manager.ActorAction {
	return map[string]entity_manager.ActorAction{
		ActionEnabled: {
			Name:        ActionEnabled,
			Description: "enable",
		},
		ActionDisable: {
			Name:        ActionDisable,
			Description: "disable",
		},
	}
}

`

var (
	controllerCmd = &cobra.Command{
		Use:   "p",
		Short: "plugin generator",
		Long:  "$ cli g p [pluginName]",
	}
	endpointName string
)

func init() {
	generate.Generate.AddCommand(controllerCmd)
	controllerCmd.Flags().StringVarP(&endpointName, "endpoint", "e", "EndpointName", "EndpointName")
	controllerCmd.Run = func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			log.Error("Wrong number of arguments. Run: cli help generate")
			return
		}

		currpath, _ := os.Getwd()

		g := Generator{}
		g.Generate(args[0], currpath)
	}
}

// Generator ...
type Generator struct{}

// Generate ...
func (e Generator) Generate(pluginName, currpath string) {

	log.Infof("Using '%s' as controller name", pluginName)

	fp := path.Join(currpath, "plugins", strings.ToLower(pluginName))

	e.addPLugin(fp, pluginName)
}

func createFile(fp, tpl, fileName, pluginName string) {

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		if err := os.MkdirAll(fp, 0777); err != nil {
			log.Errorf("Could not create plugin directory: %s", err.Error())
			return
		}
	}

	templateData := struct {
		Package    string
		PluginName string
		Dir        string
	}{
		Dir:        common.Dir(),
		Package:    strings.ToLower(pluginName),
		PluginName: pluginName,
	}

	fpath := path.Join(fp, strings.ToLower(fileName)+".go")
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		log.Errorf("Could not create file: %s", err.Error())
		return
	}
	defer f.Close()

	t := template.Must(template.New("controller").Parse(tpl))

	if t.Execute(f, templateData) != nil {
		log.Error(err.Error())
	}

	common.FormatSourceCode(fpath)
}

func (e Generator) addPLugin(fp, pluginName string) {
	createFile(fp, actorTpl, "actor", pluginName)
	createFile(fp, pluginTpl, "plugin", pluginName)
	createFile(fp, typesTpl, "types", pluginName)
}
