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

package local_migrations

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/e154/smart-home/internal/system/initial/assertions"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"
)

type MigrationTemplates struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationTemplates(adaptors *adaptors.Adaptors) *MigrationTemplates {
	return &MigrationTemplates{
		adaptors: adaptors,
	}
}

func (t *MigrationTemplates) Up(ctx context.Context) (err error) {

	dataDir := filepath.Join("data", "templates")

	var files []os.DirEntry
	files, err = os.ReadDir(dataDir)
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

		templateType := models.TemplateTypeItem
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
			templateType = models.TemplateTypeTemplate
		case "register_admin_created":
			templateType = models.TemplateTypeTemplate
		}

		var tpl *models.Template
		if templateType == models.TemplateTypeTemplate {
			tpl, err = t.adaptors.Template.GetByName(ctx, name)
		} else {
			tpl, err = t.adaptors.Template.GetItemByName(ctx, name)
		}

		if err == nil || tpl != nil {
			continue
		}

		var b []byte
		b, err = ioutil.ReadFile(filepath.Join(dataDir, fmt.Sprintf("%s.html", name)))
		So(err, ShouldBeNil)

		template := &models.Template{
			Name:       name,
			Content:    string(b),
			Status:     models.TemplateStatusActive,
			Type:       templateType,
			ParentName: parent,
		}

		err = t.adaptors.Template.Create(ctx, template)
		So(err, ShouldBeNil)
	}

	return
}
