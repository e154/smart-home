// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type AuthorizationReq struct {
	Proto   string `json:"proto,omitempty"`
	Stream  string `json:"stream,omitempty"`
	Channel string `json:"channel,omitempty"`
	Token   string `json:"token,omitempty"`
	IP      string `json:"ip,omitempty"`
}

type AuthorizationRes struct {
	Status string `json:"status,omitempty"`
}

func RemoteAuthorization(proto string, stream string, channel string, token string, ip string) bool {

	if !Storage.ServerTokenEnable() {
		return true
	}

	buf, err := json.Marshal(&AuthorizationReq{proto, stream, channel, token, ip})

	if err != nil {
		return false
	}

	request, err := http.NewRequest("POST", Storage.ServerTokenBackend(), bytes.NewBuffer(buf))

	if err != nil {
		return false
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	response, err := client.Do(request)

	if err != nil {
		return false
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return false
	}

	var res AuthorizationRes

	err = json.Unmarshal(bodyBytes, &res)

	if err != nil {
		return false
	}

	if res.Status == "1" {
		return true
	}

	return false
}
