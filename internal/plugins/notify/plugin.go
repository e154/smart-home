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

package notify

import (
	"context"
	"embed"

	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/logger"
	"github.com/e154/smart-home/pkg/plugins"
)

var (
	log = logger.MustGetLogger("plugins.notify")
)

var _ plugins.Pluggable = (*plugin)(nil)

//go:embed Readme.md
//go:embed Readme.ru.md
var F embed.FS

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
}

// New ...
func New() plugins.Pluggable {
	p := &plugin{
		Plugin: plugins.NewPlugin(),
	}
	p.F = F
	return p
}

// Load ...
func (p *plugin) Load(ctx context.Context, service plugins.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, nil); err != nil {
		return
	}

	p.Service.ScriptService().PushStruct("notifr", NewNotifyBind(service.EventBus()))
	p.Service.ScriptService().PushStruct("template", NewTemplateBind(service.Adaptors()))

	return nil
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {

	p.Service.ScriptService().PopStruct("notifr")
	p.Service.ScriptService().PopStruct("template")

	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	return nil
}

// Name ...
func (p *plugin) Name() string {
	return Name
}

// Depends ...
func (p *plugin) Depends() []string {
	return nil
}
