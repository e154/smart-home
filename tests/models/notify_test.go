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

package models

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/notify"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNotify(t *testing.T) {

	const path = "conf/notify.json"

	Convey("notify", t, func(ctx C) {
		err := container.Invoke(func(migrations *migrations.Migrations) {
			err := migrations.Purge()
			So(err, ShouldBeNil)
		})

		err = container.Invoke(func(adaptors *adaptors.Adaptors,
			notifyService *notify.Notify) {

			// read config file
			//var file []byte
			//file, err = ioutil.ReadFile(path)
			//So(err, ShouldBeNil)
			//
			//conf := &notify.NotifyConfig{}
			//err = json.Unmarshal(file, &conf)
			//So(err, ShouldBeNil)
			//
			//notifyService.UpdateCfg(conf)
			//notifyService.Restart()
			//
			//sms := notify.NewSMS()
			//sms.Phone = "+16152059974"
			//sms.Text = "test"
			//notifyService.Send(sms)
			//
			//time.Sleep(time.Second * 1)
			//
			//fmt.Println(notifyService.Stat())
			//
			//time.Sleep(time.Second * 3)
		})

		if err != nil {
			print(err.Error())
		}
	})
}
