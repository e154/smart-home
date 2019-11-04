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
