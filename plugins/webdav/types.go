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
	"fmt"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"net/http"
	"regexp"
)

const (
	// Name ...
	Name = "webdav"
	// TopicSystemStart ...
	TopicSystemStart = "system/event/start" //todo move
	// TopicSystemStop ...
	TopicSystemStop = "system/event/stop" //todo move

	AttrUser      = "user"
	AttrPassword  = "password"
	AttrAnonymous = "anonymous"

	Version = "0.0.1"
)

// NewSettings ...
func NewSettings() m.Attributes {
	return m.Attributes{
		AttrUser: {
			Name: AttrUser,
			Type: common.AttributeString,
		},
		AttrPassword: {
			Name: AttrPassword,
			Type: common.AttributeEncrypted,
		},
		AttrAnonymous: {
			Name: AttrAnonymous,
			Type: common.AttributeBool,
		},
	}
}

func scriptExt(script *m.Script) (ext string) {
	switch script.Lang {
	case common.ScriptLangTs:
		ext = "ts"
	case common.ScriptLangCoffee:
		ext = "coffee"
	case common.ScriptLangJavascript:
		ext = "js"
	default:
		ext = "txt"
	}
	return
}

func getFileName(script *m.Script) string {
	return fmt.Sprintf("%s.%s", script.Name, scriptExt(script))
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}
