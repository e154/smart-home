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

package scripts

import (
	"context"
	"fmt"
	"testing"

	"github.com/e154/smart-home/internal/system/migrations"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/scripts"

	. "github.com/smartystreets/goconvey/convey"
)

func Test9(t *testing.T) {

	//var state string
	//store = func(i interface{}) {
	//	state = fmt.Sprintf("%v", i)
	//}

	const path = "conf/notify.json"

	Convey("clear db", t, func(ctx C) {
		err := container.Invoke(func(migrations *migrations.Migrations) {

			// clear database
			// ------------------------------------------------
			err := migrations.Purge()
			So(err, ShouldBeNil)
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	Convey("send sms", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService) {

			// add templates
			// ------------------------------------------------
			templates := []*models.Template{
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
			}

			for _, template := range templates {
				err := adaptors.Template.UpdateOrCreate(context.Background(), template)
				So(err, ShouldBeNil)
			}

			// ------------------------------------------------
			render, err := adaptors.Template.Render(context.Background(), "template2", map[string]interface{}{
				"code": 12345,
			})
			So(err, ShouldBeNil)
			So(render.Subject, ShouldEqual, "")
			So(render.Body, ShouldEqual, "Activate code: 12345")

			//read config file
			// ------------------------------------------------
			//var file []byte
			//file, err = ioutil.ReadFile(path)
			//So(err, ShouldBeNil)
			//
			//conf := &notify.Config{}
			//err = json.Unmarshal(file, &conf)
			//So(err, ShouldBeNil)
			//
			//notifyService.UpdateCfg(conf)
			//notifyService.Restart()
			//
			//// scripts
			//// ------------------------------------------------
			//storeRegisterCallback(scriptService)
			//
			//scripts := GetScripts(ctx, scriptService, adaptors, 24)
			//
			//engine24, err := scriptService.NewEngine(scripts["script24"])
			//So(err, ShouldBeNil)
			//err = engine24.Compile()
			//So(err, ShouldBeNil)
			//
			//_, err = engine24.DoCustom("main")
			//So(err, ShouldBeNil)
			//
			//time.Sleep(time.Second * 5)
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	})
}
