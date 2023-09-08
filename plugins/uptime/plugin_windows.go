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

package uptime

import (
	"time"

	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

const (
	name = "uptime"
)

// plugin ...
type plugin struct {
	*supervisor.Plugin
	entity     *Actor
	ticker     *time.Ticker
	pause      time.Duration
	storyModel *m.RunStory
	quit       chan struct{}
}

// Load ...
func (p *plugin) Load(service supervisor.Service) (err error) {

	return
}

// Unload ...
func (u *plugin) Unload() (err error) {

	return
}

// Name ...
func (u plugin) Name() string {
	return name
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
	return "0.0.1"
}
