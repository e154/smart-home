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

package endpoint

// PluginActorEndpoint ...
type PluginActorEndpoint struct {
	*CommonEndpoint
}

// NewPluginActorEndpoint ...
func NewPluginActorEndpoint(common *CommonEndpoint) *PluginActorEndpoint {
	return &PluginActorEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (p *PluginActorEndpoint) Add() {}

// GetByName ...
func (p *PluginActorEndpoint) GetByName(name string) {}

// Update ...
func (p *PluginActorEndpoint) Update() {}

// Delete ...
func (p *PluginActorEndpoint) Delete() {}

// Search ...
func (p *PluginActorEndpoint) Search(query string, limit, offset int) {}
