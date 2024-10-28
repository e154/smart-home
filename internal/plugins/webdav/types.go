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
	"path/filepath"
	"strings"
	"time"

	commonPkg "github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
)

const (
	// Name ...
	Name = "webdav"

	AttrUser      = "user"
	AttrPassword  = "password"
	AttrAnonymous = "anonymous"
)

// NewSettings ...
func NewSettings() m.Attributes {
	return m.Attributes{
		AttrUser: {
			Name: AttrUser,
			Type: commonPkg.AttributeString,
		},
		AttrPassword: {
			Name: AttrPassword,
			Type: commonPkg.AttributeEncrypted,
		},
		AttrAnonymous: {
			Name: AttrAnonymous,
			Type: commonPkg.AttributeBool,
		},
	}
}

func extractScriptName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.ReplaceAll(base, ext, "")
}

func getScriptLang(path string) commonPkg.ScriptLang {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	switch ext {
	case ".ts":
		return "ts"
	case ".js":
		return "javascript"
	case ".coffee":
		return "coffeescript"
	}
	return ""
}

func getExtension(script *m.Script) (ext string) {
	switch script.Lang {
	case commonPkg.ScriptLangTs:
		ext = ".ts"
	case commonPkg.ScriptLangCoffee:
		ext = ".coffee"
	case commonPkg.ScriptLangJavascript:
		ext = ".js"
	default:
		ext = ".txt"
	}
	return
}

func getFileName(script *m.Script) string {
	return fmt.Sprintf("%s%s", script.Name, getExtension(script))
}

type FileInfo struct {
	Size          int64
	ModTime       time.Time
	LastCheck     time.Time
	IsInitialized bool
}
