// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package models

import (
	"fmt"
	"testing"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/endpoint"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/migrations"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTemplate(t *testing.T) {

	const subject = "Lorem ipsum dolor sit amet"
	const body = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

	Convey("add user", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			endpoint *endpoint.Endpoint) {

			// clear database
			// ------------------------------------------------
			err := migrations.Purge()
			So(err, ShouldBeNil)

			// add templates
			// ------------------------------------------------
			templates := []*m.Template{
				{
					Name:       "main",
					Content:    "[message:block]",
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeItem,
					ParentName: nil,
				},
				{
					Name:       "message",
					Content:    "[title:block][body:block]",
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeItem,
					ParentName: common.String("main"),
				},
				{
					Name:       "title",
					Content:    "[title:content] [var1]",
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeItem,
					ParentName: common.String("message"),
				},
				{
					Name:       "body",
					Content:    "[body:content] [var2]",
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeItem,
					ParentName: common.String("message"),
				},
				{
					Name: "template1",
					Content: fmt.Sprintf(`{
  "items": [
    "title",
    "body"
  ],
  "title": "%s",
  "fields": [
    {
      "name": "title",
      "value": "Lorem ipsum dolor sit amet,"
    },
    {
      "name": "body",
      "value": ", sed do eiusmod tempor incididunt"
    }
  ]
}`, subject),
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeTemplate,
					ParentName: nil,
				},
				{
					Name:       "sms_body",
					Content:    "[code:block]",
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeItem,
					ParentName: nil,
				},
				{
					Name:       "code",
					Content:    "[code:content] [code]",
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeItem,
					ParentName: common.String("sms_body"),
				},
				{
					Name: "template2",
					Content: fmt.Sprintf(`{
 "items": [
   "code"
 ],
 "title": "",
 "fields": [
	{
     "name": "code",
     "value": "Activate code:"
   }
]
}`),
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeTemplate,
					ParentName: nil,
				},
				{
					Name:       "sms_warning",
					Content:    "some warning message",
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeItem,
					ParentName: nil,
				},
				{
					Name: "template3",
					Content: fmt.Sprintf(`{
 "items": [
   "sms_warning"
 ],
 "title": "",
 "fields": []
}`),
					Status:     m.TemplateStatusActive,
					Type:       m.TemplateTypeTemplate,
					ParentName: nil,
				},
			}

			for _, template := range templates {
				err := adaptors.Template.UpdateOrCreate(template)
				So(err, ShouldBeNil)
			}

			// Lorem ipsum dolor sit amet
			// ------------------------------------------------
			render, err := adaptors.Template.Render("template1", map[string]interface{}{
				"var1": "consectetur adipiscing elit",
				"var2": "ut labore et dolore magna aliqua.",
			})
			So(err, ShouldBeNil)
			So(render.Subject, ShouldEqual, subject)
			So(render.Body, ShouldEqual, body)

			// Activate code: 12345
			// ------------------------------------------------
			render, err = adaptors.Template.Render("template2", map[string]interface{}{
				"code": 12345,
			})
			So(err, ShouldBeNil)
			So(render.Subject, ShouldEqual, "")
			So(render.Body, ShouldEqual, "Activate code: 12345")

			// warning message
			// ------------------------------------------------
			render, err = adaptors.Template.Render("template3", nil)
			So(err, ShouldBeNil)
			So(render.Subject, ShouldEqual, "")
			So(render.Body, ShouldEqual, "some warning message")

			//...
		})
	})
}
