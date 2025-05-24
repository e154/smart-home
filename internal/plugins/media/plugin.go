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

package media

import (
	"context"
	"embed"
	"net/http"

	"github.com/e154/smart-home/internal/plugins/media/server"
	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/logger"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

var (
	log = logger.MustGetLogger("plugins.media")
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
	server *server.Media
	router *http.ServeMux
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

	p.server = server.NewMedia(service.EventBus())
	if err = p.server.Start(); err != nil {
		return
	}

	controller := NewControllerMedia()
	p.router = http.NewServeMux()

	//a.echo.Any("/stream/:entity_id/channel/:channel/mse", a.echoFilter.Auth(a.controllers.StreamMSE)) //Auth
	//a.echo.Any("/stream/:entity_id/channel/:channel/hlsll/live/init.mp4", a.controllers.Media.StreamHLSLLInit)
	//a.echo.Any("/stream/:entity_id/channel/:channel/hlsll/live/index.m3u8", a.controllers.Media.StreamHLSLLM3U8)
	//a.echo.Any("/stream/:entity_id/channel/:channel/hlsll/live/segment/:segment/:any", a.controllers.Media.StreamHLSLLM4Segment)
	//a.echo.Any("/stream/:entity_id/channel/:channel/hlsll/live/fragment/:segment/:fragment/:any", a.controllers.Media.StreamHLSLLM4Fragment)

	p.router.Handle("/media/{entity_id}/channel/{channel}/mse", p.Service.HttpAccessFilter().Auth(http.HandlerFunc(controller.StreamMSE)))

	return nil
}

func (p *plugin) Unload(ctx context.Context) error {

	p.server.Shutdown()
	_ = p.Plugin.Unload(ctx)
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

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{}
}

// ServeHTTP ...
func (p *plugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}
