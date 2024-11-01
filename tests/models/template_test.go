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

package models

import (
	"context"
	"fmt"
	"testing"

	"github.com/e154/smart-home/internal/endpoint"
	"github.com/e154/smart-home/internal/system/migrations"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"

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
			templates := []*models.Template{
				{
					Name:       "main",
					Content:    "[message:block]",
					Status:     models.TemplateStatusActive,
					Type:       models.TemplateTypeItem,
					ParentName: nil,
				},
				{
					Name:       "message",
					Content:    "[title:block][body:block]",
					Status:     models.TemplateStatusActive,
					Type:       models.TemplateTypeItem,
					ParentName: common.String("main"),
				},
				{
					Name:       "title",
					Content:    "[title:content] [var1]",
					Status:     models.TemplateStatusActive,
					Type:       models.TemplateTypeItem,
					ParentName: common.String("message"),
				},
				{
					Name:       "body",
					Content:    "[body:content] [var2]",
					Status:     models.TemplateStatusActive,
					Type:       models.TemplateTypeItem,
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
					Status:     models.TemplateStatusActive,
					Type:       models.TemplateTypeTemplate,
					ParentName: nil,
				},
				{
					Name:       "sms_body",
					Content:    "[code:block]",
					Status:     models.TemplateStatusActive,
					Type:       models.TemplateTypeItem,
					ParentName: nil,
				},
				{
					Name:       "code",
					Content:    "[code:content] [code]",
					Status:     models.TemplateStatusActive,
					Type:       models.TemplateTypeItem,
					ParentName: common.String("sms_body"),
				},
				{
					Name: "template2",
					Content: `{
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
}`,
					Status:     models.TemplateStatusActive,
					Type:       models.TemplateTypeTemplate,
					ParentName: nil,
				},
				{
					Name:       "sms_warning",
					Content:    "some warning message",
					Status:     models.TemplateStatusActive,
					Type:       models.TemplateTypeItem,
					ParentName: nil,
				},
				{
					Name: "template3",
					Content: `{
 "items": [
   "sms_warning"
 ],
 "title": "",
 "fields": []
}`,
					Status:     models.TemplateStatusActive,
					Type:       models.TemplateTypeTemplate,
					ParentName: nil,
				},
			}

			for _, template := range templates {
				err := adaptors.Template.UpdateOrCreate(context.Background(), template)
				So(err, ShouldBeNil)
			}

			// Lorem ipsum dolor sit amet
			// ------------------------------------------------
			render, err := adaptors.Template.Render(context.Background(), "template1", map[string]interface{}{
				"var1": "consectetur adipiscing elit",
				"var2": "ut labore et dolore magna aliqua.",
			})
			So(err, ShouldBeNil)
			So(render.Subject, ShouldEqual, subject)
			So(render.Body, ShouldEqual, body)

			// Activate code: 12345
			// ------------------------------------------------
			render, err = adaptors.Template.Render(context.Background(), "template2", map[string]interface{}{
				"code": 12345,
			})
			So(err, ShouldBeNil)
			So(render.Subject, ShouldEqual, "")
			So(render.Body, ShouldEqual, "Activate code: 12345")

			// warning message
			// ------------------------------------------------
			render, err = adaptors.Template.Render(context.Background(), "template3", nil)
			So(err, ShouldBeNil)
			So(render.Subject, ShouldEqual, "")
			So(render.Body, ShouldEqual, "some warning message")

			//...
		})
	})
}
