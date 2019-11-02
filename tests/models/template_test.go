package models

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/endpoint"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/migrations"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestTemplate(t *testing.T) {

	Convey("add user", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			endpoint *endpoint.Endpoint) {

			//count, items, err := endpoint.Template.GetItemsSortedList()
			//So(err, ShouldBeNil)
			//
			//fmt.Println(count)
			//fmt.Println(items)

			dataDir := filepath.Join("data", "templates")

			files, err := ioutil.ReadDir(dataDir)
			So(err, ShouldBeNil)

			for _, file := range files {
				fmt.Println(file)
				if file.IsDir() {
					continue
				}

				b, err := ioutil.ReadFile(filepath.Join(dataDir, file.Name()))
				So(err, ShouldBeNil)

				templateType := m.TemplateTypeItem
				var parent *string

				name := strings.Replace(file.Name(), ".html", "", -1)

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
					parent = common.String("social")
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

				template := &m.Template{
					Name:       name,
					Content:    string(b),
					Status:     m.TemplateStatusActive,
					Type:       templateType,
					ParentName: parent,
				}

				err = adaptors.Template.UpdateOrCreate(template)
				So(err, ShouldBeNil)
				return
			}
		})
	})
}
