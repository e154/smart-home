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

package models

// Plugin ...
type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Enabled bool   `json:"enabled"`
	System  bool   `json:"system"`
}

type PluginOptions struct {
	Triggers           bool                         `json:"triggers"`
	Actors             bool                         `json:"actors"`
	ActorCustomAttrs   bool                         `json:"actor_custom_attrs"`
	ActorAttrs         Attributes                   `json:"actor_attrs"`
	ActorCustomActions bool                         `json:"actor_custom_actions"`
	ActorActions       map[string]EntityActionShort `json:"actor_actions"`
	ActorCustomStates  bool                         `json:"actor_custom_states"`
	ActorStates        map[string]EntityStateShort  `json:"actor_states"`
	ActorSetts         Attributes                   `json:"actor_setts"`
}
