package models

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/migrations"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTemplate(t *testing.T) {

	Convey("add user", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			endpoint *endpoint.Endpoint) {

			count, items, err := endpoint.Template.GetItemsSortedList()
			So(err, ShouldBeNil)

			fmt.Println(count)
			fmt.Println(items)
		})
	})
}
