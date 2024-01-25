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

package webdav

import (
	"context"
	"embed"
	"net/http"

	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.webdav")
)

var _ supervisor.Pluggable = (*plugin)(nil)

//go:embed Readme.md
//go:embed Readme.ru.md
var F embed.FS

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	server   *Server
	settings m.Attributes
}

// New ...
func New() supervisor.Pluggable {
	p := &plugin{
		Plugin: supervisor.NewPlugin(),
	}
	p.F = F
	return p
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, nil); err != nil {
		return
	}

	// load settings
	p.settings, err = p.LoadSettings(p)
	if err != nil {
		log.Warn(err.Error())
		p.settings = NewSettings()
	}

	if p.settings == nil {
		p.settings = NewSettings()
	}

	p.server = NewServer()
	p.server.Start(service.Adaptors(), service.ScriptService(), service.EventBus())

	return
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	p.server.Shutdown()

	return
}

// Name ...
func (p plugin) Name() string {
	return Name
}

// Type ...
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginBuiltIn
}

// Depends ...
func (p *plugin) Depends() []string {
	return nil
}

// Version ...
func (p *plugin) Version() string {
	return Version
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Javascript: m.PluginOptionsJs{},
		Setts:      NewSettings(),
	}
}

// ServeHTTP ...
func (p *plugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()

	if p.settings[AttrAnonymous].Bool() ||
		username == p.settings[AttrUser].String() && password == p.settings[AttrPassword].Decrypt() {
		p.server.ServeHTTP(w, r)
		return
	}

	w.Header().Set("WWW-Authenticate", `Basic realm="BASIC WebDAV REALM"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
}
