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

package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/e154/smart-home/version"
)

// ControllerIndex ...
type ControllerIndex struct {
	*ControllerCommon
}

// NewControllerIndex ...
func NewControllerIndex(common *ControllerCommon) *ControllerIndex {
	return &ControllerIndex{
		ControllerCommon: common,
	}
}

// Index ...
func (c ControllerIndex) Index(publicAssets fs.FS) http.Handler {
	serverVersion, _ := json.Marshal(version.GetVersion())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := map[string]interface{}{
			"server_url":     c.ControllerCommon.appConfig.ApiFullAddress(),
			"run_mode":       c.ControllerCommon.appConfig.Mode,
			"server_version": string(serverVersion),
		}
		templates := template.Must(template.New("index").ParseFS(publicAssets, "index.html"))

		err := templates.ExecuteTemplate(w, "index.html", &b)
		if err != nil {
			http.Error(w, fmt.Sprintf("index: couldn't parse template: %v", err), http.StatusInternalServerError)
			return
		}
	})
}
