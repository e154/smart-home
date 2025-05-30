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

package system

import (
	"context"
	"embed"
	"sync"

	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/logger"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/plugins/triggers"
)

var (
	log = logger.MustGetLogger("plugins.system")
)

var _ plugins.Pluggable = (*plugin)(nil)

//go:embed *.md
var F embed.FS

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	actorsLock *sync.Mutex
	registrar  triggers.IRegistrar
}

// New ...
func New() plugins.Pluggable {
	p := &plugin{
		Plugin:     plugins.NewPlugin(),
		actorsLock: &sync.Mutex{},
	}
	p.F = F
	return p
}

// Load ...
func (p *plugin) Load(ctx context.Context, service plugins.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, nil); err != nil {
		return
	}

	// register trigger
	if triggersPlugin, ok := service.Plugins()[triggers.Name]; ok {
		if p.registrar, ok = triggersPlugin.(triggers.IRegistrar); ok {
			if err = p.registrar.RegisterTrigger(NewTrigger(p.Service.EventBus())); err != nil {
				log.Error(err.Error())
				return
			}
		}
	}

	return nil
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	if err = p.registrar.UnregisterTrigger(Name); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// Name ...
func (p *plugin) Name() string {
	return Name
}

// Depends ...
func (p *plugin) Depends() []string {
	return []string{"triggers"}
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Triggers:      true,
		TriggerParams: NewTriggerParams(),
	}
}
