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

package system

import (
	"fmt"
	"github.com/e154/smart-home/system/message_queue"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
	"time"
)

func TestMessageQueue(t *testing.T) {

	const (
		queueSize = 10
	)

	t.Run("base topic", func(t *testing.T) {
		Convey("case", t, func(ctx C) {

			queue := message_queue.New(queueSize)

			wg := &sync.WaitGroup{}
			wg.Add(1)
			arr := make([]interface{}, 0)
			fn := func(topic string, args ...interface{}) {
				arr = append(arr, args...)
				wg.Done()
			}
			queue.Subscribe("a/b", fn)

			queue.Publish("a", "msg1", "msg1")
			queue.Publish("a/b", "msg2", "msg2")
			queue.Publish("a/b/c", "msg3", "msg3")
			queue.Publish("a/b/c/d", "msg4", "msg4")

			wg.Wait()

			queue.Unsubscribe("a/b", fn)

			ctx.So(fmt.Sprintf("%v", arr), ShouldEqual, "[msg2 msg2]")
		})
	})

	t.Run("multi level", func(t *testing.T) {
		Convey("case", t, func(ctx C) {

			queue := message_queue.New(queueSize)

			wg := &sync.WaitGroup{}
			wg.Add(3)
			arr := make([]string, 0)
			fn := func(topic string, msg string) {
				arr = append(arr, msg)
				wg.Done()
			}
			queue.Subscribe("a/b/#", fn)

			queue.Publish("a", "msg1")
			queue.Publish("a/b", "msg2")
			queue.Publish("a/b/c", "msg3")
			queue.Publish("a/b/c/d", "msg4")

			wg.Wait()

			queue.Unsubscribe("a/b/#", fn)

			ctx.So(fmt.Sprintf("%v", arr), ShouldEqual, "[msg2 msg3 msg4]")
		})
	})

	t.Run("single level", func(t *testing.T) {
		Convey("case", t, func(ctx C) {

			queue := message_queue.New(queueSize)

			wg := &sync.WaitGroup{}
			wg.Add(1)
			arr := make([]string, 0)
			fn := func(topic string, msg string) {
				arr = append(arr, msg)
				wg.Done()
			}
			queue.Subscribe("a/b/+/d", fn)

			queue.Publish("a", "msg1")
			queue.Publish("a/b", "msg2")
			queue.Publish("a/b/c", "msg3")
			queue.Publish("a/b/c/d", "msg4")
			queue.Publish("a/b/c/d/e", "msg5")
			queue.Publish("a/b/c/d/e/f", "msg6")

			wg.Wait()

			queue.Unsubscribe("a/b/+/d", fn)

			ctx.So(fmt.Sprintf("%v", arr), ShouldEqual, "[msg4]")
		})
	})

	t.Run("subscribing", func(t *testing.T) {
		Convey("case", t, func(ctx C) {

			queue := message_queue.New(queueSize)

			arr := make([]string, 0)
			fn := func(topic string, msg ...string) {
				arr = append(arr, msg...)
			}
			queue.Subscribe("a/#", fn)

			queue.Publish("a", "msg1")
			queue.Publish("a/b", "msg2")

			time.Sleep(time.Millisecond * 500)

			queue.Unsubscribe("a/#", fn)

			queue.Publish("a", "msg3")
			queue.Publish("a/b", "msg4")

			time.Sleep(time.Millisecond * 500)

			ctx.So(fmt.Sprintf("%v", arr), ShouldEqual, "[msg1 msg2]")
		})
	})
}
