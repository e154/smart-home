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

package api

import (
	"context"
	"fmt"
	"testing"

	gw "github.com/e154/smart-home/api/stub/api"
	m "github.com/e154/smart-home/models"
	container2 "github.com/e154/smart-home/tests/api/container"
	"google.golang.org/grpc"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/migrations"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVariables(t *testing.T) {

	Convey("variables", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			dialer *container2.Dialer) {

			err := migrations.Purge()
			ctx.So(err, ShouldBeNil)

			v1 := m.Variable{
				Name:  "v1",
				Value: "v1",
			}
			err = adaptors.Variable.Add(context.Background(), v1)
			ctx.So(err, ShouldBeNil)

			c := context.Background()
			conn, err := grpc.DialContext(c, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer.Call()))
			ctx.So(err, ShouldBeNil)
			defer conn.Close()
			client := gw.NewVariableServiceClient(conn)

			t.Run("add", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					v2, err := client.AddVariable(c, &gw.NewVariableRequest{
						Name:  "v2",
						Value: "v2",
					})
					ctx.So(err, ShouldBeNil)
					ctx.So(v2.Name, ShouldEqual, "v2")
					ctx.So(v2.Value, ShouldEqual, "v2")

					v1, err := client.AddVariable(c, &gw.NewVariableRequest{
						Name:  "v1",
						Value: "v11",
					})
					ctx.So(err, ShouldBeNil)
					ctx.So(v1.Name, ShouldEqual, "v1")
					ctx.So(v1.Value, ShouldEqual, "v11")

					v3, err := client.AddVariable(c, &gw.NewVariableRequest{
						Name:  "v3",
						Value: "v3",
					})
					ctx.So(err, ShouldBeNil)
					ctx.So(v3.Name, ShouldEqual, "v3")
					ctx.So(v3.Value, ShouldEqual, "v3")

				})
			})

			t.Run("getByName", func(t *testing.T) {
				Convey("", t, func(ctx C) {
					v1, err := client.GetVariableByName(c, &gw.GetVariableRequest{
						Name: "v1",
					})
					ctx.So(err, ShouldBeNil)
					ctx.So(v1.Name, ShouldEqual, "v1")
					ctx.So(v1.Value, ShouldEqual, "v11")

					v3, err := client.GetVariableByName(c, &gw.GetVariableRequest{
						Name: "v3",
					})
					ctx.So(err, ShouldBeNil)
					ctx.So(v3.Name, ShouldEqual, "v3")
					ctx.So(v3.Value, ShouldEqual, "v3")

					_, err = client.GetVariableByName(c, &gw.GetVariableRequest{
						Name: "v4",
					})
					ctx.So(err, ShouldNotBeNil)
				})
			})

			t.Run("update", func(t *testing.T) {
				Convey("", t, func(ctx C) {

					v3, err := client.UpdateVariable(c, &gw.UpdateVariableRequest{
						Name:  "v3",
						Value: "v333",
					})
					ctx.So(err, ShouldBeNil)
					ctx.So(v3.Name, ShouldEqual, "v3")
					ctx.So(v3.Value, ShouldEqual, "v333")

					v3, err = client.GetVariableByName(c, &gw.GetVariableRequest{
						Name: "v3",
					})
					ctx.So(err, ShouldBeNil)
					ctx.So(v3.Name, ShouldEqual, "v3")
					ctx.So(v3.Value, ShouldEqual, "v333")

					_, err = client.GetVariableByName(c, &gw.GetVariableRequest{
						Name: "v4",
					})
					ctx.So(err, ShouldNotBeNil)
				})
			})

			t.Run("list", func(t *testing.T) {
				Convey("", t, func(ctx C) {
					list, err := client.GetVariableList(c, &gw.PaginationRequest{})
					ctx.So(err, ShouldBeNil)
					ctx.So(len(list.Items), ShouldEqual, 3)
					ctx.So(list.Meta.Total, ShouldEqual, 3)
				})
			})

			t.Run("delete", func(t *testing.T) {
				Convey("", t, func(ctx C) {
					_, err = client.DeleteVariable(c, &gw.DeleteVariableRequest{Name: "v4"})
					ctx.So(err, ShouldBeNil)

					_, err := client.DeleteVariable(c, &gw.DeleteVariableRequest{Name: "v3"})
					ctx.So(err, ShouldBeNil)

					list, err := client.GetVariableList(c, &gw.PaginationRequest{})
					ctx.So(err, ShouldBeNil)
					ctx.So(len(list.Items), ShouldEqual, 2)
					ctx.So(list.Meta.Total, ShouldEqual, 2)
				})
			})
		})

		if err != nil {
			fmt.Println(err.Error())
		}
	})
}
