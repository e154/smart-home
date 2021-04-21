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

package env1

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// TemplateManager ...
type TemplateManager struct {
	adaptors *adaptors.Adaptors
}

// NewTemplateManager ...
func NewTemplateManager(adaptors *adaptors.Adaptors) *TemplateManager {
	return &TemplateManager{
		adaptors: adaptors,
	}
}

// Create ...
func (t TemplateManager) Create() {

	dataDir := filepath.Join("data", "templates")

	files, err := os.ReadDir(dataDir)
	So(err, ShouldBeNil)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := strings.Replace(file.Name(), ".html", "", -1)

		switch name {
		case "main", "message", "body", "callout", "footer", "contacts", "social", "facebook", "google", "header", "password_reset", "privacy", "register_admin_created", "title", "twitter", "vk":
			continue
		default:
			//fmt.Printf("unknown file %v", file.Name())
			return
		}
	}

	fileNames := []string{"main", "message", "body", "callout", "footer", "contacts", "social", "facebook", "google", "header", "password_reset", "privacy", "register_admin_created", "title", "twitter", "vk"}

	for _, name := range fileNames {

		templateType := m.TemplateTypeItem
		var parent *string

		switch name {
		case "header":
			parent = common.String("main")
		case "main":
		case "body":
			parent = common.String("message")
		case "google":
			parent = common.String("social")
		case "message":
			parent = common.String("main")
		case "callout":
			parent = common.String("main")
		case "contacts":
			parent = common.String("footer")
		case "footer":
			parent = common.String("main")
		case "social":
			parent = common.String("footer")
		case "facebook":
			parent = common.String("social")
		case "privacy":
			parent = common.String("main")
		case "title":
			parent = common.String("message")
		case "twitter":
			parent = common.String("social")
		case "vk":
			parent = common.String("social")
		case "password_reset":
			templateType = m.TemplateTypeTemplate
		case "register_admin_created":
			templateType = m.TemplateTypeTemplate
		}

		var tpl *m.Template
		if templateType == m.TemplateTypeTemplate {
			tpl, err = t.adaptors.Template.GetByName(name)
		} else {
			tpl, err = t.adaptors.Template.GetItemByName(name)
		}

		if err == nil || tpl != nil {
			continue
		}

		b, err := ioutil.ReadFile(filepath.Join(dataDir, fmt.Sprintf("%s.html", name)))
		So(err, ShouldBeNil)

		template := &m.Template{
			Name:       name,
			Content:    string(b),
			Status:     m.TemplateStatusActive,
			Type:       templateType,
			ParentName: parent,
		}

		err = t.adaptors.Template.Create(template)
		So(err, ShouldBeNil)
	}
}

// Upgrade ...
func (t TemplateManager) Upgrade(oldVersion int) (err error) {

	switch oldVersion {
	case 0:

	}

	return
}
