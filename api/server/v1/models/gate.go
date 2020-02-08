// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

// swagger:model
type GateSettings struct {
	GateServerToken string `json:"gate_server_token"`
	Address         string `json:"address"`
	Enabled         bool   `json:"enabled"`
}

// swagger:model
type UpdateGateSettings struct {
	GateServerToken string `json:"gate_server_token"`
	Address         string `json:"address"`
	Enabled         bool   `json:"enabled"`
}

// swagger:model
type GateMobileList struct {
	Total     int64    `json:"total"`
	TokenList []string `json:"token_list"`
}

// swagger:model
type DeleteGateMobile struct {
	Token string `json:"token"`
}
