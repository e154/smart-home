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

package gate_client

import "net/http"

// Settings ...
type Settings struct {
	GateServerToken string `json:"gate_server_token"`
	Address         string `json:"address"`
	Enabled         bool   `json:"enabled"`
}

// Valid ...
func (s Settings) Valid() bool {
	if s.Address != "" && s.Enabled {
		return true
	}
	return false
}

// Equal ...
func (s Settings) Equal(v Settings) bool {
	return s.GateServerToken == v.GateServerToken &&
		s.Address == v.Address &&
		s.Enabled == v.Enabled
}

// MobileList ...
type MobileList struct {
	Total     int64    `json:"total"`
	TokenList []string `json:"token_list"`
}

// StreamRequestModel ...
type StreamRequestModel struct {
	URI    string      `json:"uri"`
	Method string      `json:"method"`
	Body   []byte      `json:"body"`
	Header http.Header `json:"header"`
}

// StreamResponseModel ...
type StreamResponseModel struct {
	Code   int         `json:"code"`
	Body   []byte      `json:"body"`
	Header http.Header `json:"header"`
}
