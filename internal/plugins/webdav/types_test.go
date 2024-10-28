// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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
	"testing"

	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"

	"github.com/stretchr/testify/require"
)

func TestExtractScriptName(t *testing.T) {

	name := extractScriptName("/webdav/scripts/foo.ts")
	require.Equal(t, "foo", name)

	name = extractScriptName("/webdav/scripts/bar.js")
	require.Equal(t, "bar", name)

	name = extractScriptName("/webdav/scripts/foo.coffee")
	require.Equal(t, "foo", name)

	name = extractScriptName("/webdav/scripts/lib.d.ts")
	require.Equal(t, "lib.d", name)
}

func TestGetScriptLang(t *testing.T) {

	lang := getScriptLang("/webdav/scripts/foo.ts")
	require.Equal(t, common.ScriptLangTs, lang)

	lang = getScriptLang("/webdav/scripts/bar.js")
	require.Equal(t, common.ScriptLangJavascript, lang)

	lang = getScriptLang("/webdav/scripts/foo.coffee")
	require.Equal(t, common.ScriptLangCoffee, lang)

	lang = getScriptLang("/webdav/scripts/lib.d.ts")
	require.Equal(t, common.ScriptLangTs, lang)
}

func TestGetExtension(t *testing.T) {

	ext := getExtension(&m.Script{
		Lang: common.ScriptLangJavascript,
	})
	require.Equal(t, ".js", ext)

	ext = getExtension(&m.Script{
		Lang: common.ScriptLangCoffee,
	})
	require.Equal(t, ".coffee", ext)

	ext = getExtension(&m.Script{
		Lang: common.ScriptLangTs,
	})
	require.Equal(t, ".ts", ext)

	ext = getExtension(&m.Script{})
	require.Equal(t, ".txt", ext)
}

func TestGetFileName(t *testing.T) {

	name := getFileName(&m.Script{
		Name: "foo",
		Lang: common.ScriptLangJavascript,
	})
	require.Equal(t, "foo.js", name)

	name = getFileName(&m.Script{
		Name: "bar",
		Lang: common.ScriptLangCoffee,
	})
	require.Equal(t, "bar.coffee", name)

	name = getFileName(&m.Script{
		Name: "foo",
		Lang: common.ScriptLangTs,
	})
	require.Equal(t, "foo.ts", name)

	name = getFileName(&m.Script{
		Name: "bar",
	})
	require.Equal(t, "bar.txt", name)
}
func TestGetFilePath(t *testing.T) {

	service := NewScripts(nil)

	path := service.getFilePath(&m.Script{
		Name: "foo",
		Lang: common.ScriptLangJavascript,
	})
	require.Equal(t, "/webdav/scripts/foo.js", path)

	path = service.getFilePath(&m.Script{
		Name: "bar",
		Lang: common.ScriptLangCoffee,
	})
	require.Equal(t, "/webdav/scripts/bar.coffee", path)

	path = service.getFilePath(&m.Script{
		Name: "foo",
		Lang: common.ScriptLangTs,
	})
	require.Equal(t, "/webdav/scripts/foo.ts", path)

	path = service.getFilePath(&m.Script{
		Name: "bar",
	})
	require.Equal(t, "/webdav/scripts/bar.txt", path)

}
