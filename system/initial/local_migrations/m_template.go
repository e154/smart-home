package local_migrations

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
)

type MigrationTemplates struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationTemplates(adaptors *adaptors.Adaptors) *MigrationTemplates {
	return &MigrationTemplates{
		adaptors: adaptors,
	}
}

func (t *MigrationTemplates) Up(ctx context.Context, adaptors *adaptors.Adaptors) (err error) {

	if adaptors != nil {
		t.adaptors = adaptors
	}
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

		template := &m.Template{
			Name:       name,
			Content:    string(b),
			Status:     m.TemplateStatusActive,
			Type:       templateType,
			ParentName: parent,
		}

		err = t.adaptors.Template.Create(ctx, template)
		So(err, ShouldBeNil)
	}

	return
}
