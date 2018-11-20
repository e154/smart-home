package workflow

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
)

// create flow
//
//  workflow + workflow scenario
//
//  emitter --> handler
//
func Test3(t *testing.T) {

	Convey("add scripts", t, func(ctx C) {
		container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			c *core.Core) {

			// clear database
			migrations.Purge()


		})
	})
}
