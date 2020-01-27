package api

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/server"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestDevice(t *testing.T) {

	var accessToken string

	type newDeviceRequest struct {
		ResponseCode int
		Device       models.NewDevice
		Id           int64
	}

	devices := []newDeviceRequest{
		{
			ResponseCode: 200,
			Device: models.NewDevice{
				Name: "device2", Description: "device desc", Status: "enabled", Type: "modbus_rtu",
				Properties: models.DeviceProperties{
					DevModBusRtuConfig: &models.DevModBusRtuConfig{
						SlaveId: 1, Baud: 115200, DataBits: 8, StopBits: 2, Parity: "none", Timeout: 100,
					},
				},
			},
		},
		{
			ResponseCode: 200,
			Device: models.NewDevice{
				Name: "device3", Description: "device desc", Status: "disabled", Type: "modbus_rtu",
				Properties: models.DeviceProperties{
					DevModBusRtuConfig: &models.DevModBusRtuConfig{
						SlaveId: 1, Baud: 115200, DataBits: 8, StopBits: 2, Parity: "none", Timeout: 100,
					},
				},
			},
		},
		{
			ResponseCode: 500,
			Device: models.NewDevice{
				Name: "device4", Description: "device desc", Status: "active", Type: "modbus_rtu",
				Properties: models.DeviceProperties{
					DevModBusRtuConfig: &models.DevModBusRtuConfig{
						SlaveId: 1, Baud: 115200, DataBits: 8, StopBits: 2, Parity: "none", Timeout: 100,
					},
				},
			},
		},
		{
			ResponseCode: 200,
			Device: models.NewDevice{
				Name: "device5", Description: "device desc", Status: "enabled", Type: "modbus_tcp",
				Properties: models.DeviceProperties{
					DevModBusTcpConfig: &models.DevModBusTcpConfig{
						SlaveId: 1, AddressPort: "127.0.0.1:502",
					},
				},
			},
		},
		{
			ResponseCode: 200,
			Device: models.NewDevice{
				Name: "device6", Description: "device desc", Status: "enabled", Type: "custom",
			},
		},
		{
			ResponseCode: 200,
			Device: models.NewDevice{
				Name: "device6", Description: "device desc", Status: "enabled", Type: "custom",
			},
		},
		{
			ResponseCode: 200,
			Device: models.NewDevice{
				Name: "device7", Description: "device desc", Status: "enabled", Type: "mqtt",
				Properties: models.DeviceProperties{
					DevMqttConfig: &models.DevMqttConfig{
						Address: "127.0.0.1", User: "user", Password: "pass",
					},
				},
			},
		},
	}

	Convey("POST /device", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core,
			accessList *access_list.AccessListService, ) {

			// clear database
			err := migrations.Purge()
			ctx.So(err, ShouldBeNil)

			// add roles
			AddRoles(adaptors, accessList, ctx)

			go server.Start()

			time.Sleep(time.Second * 1)

			client := NewClient(server.GetEngine())

			// new device params
			deviceDefault := models.NewDevice{
				Name:        "device1",
				Description: "device desc",
				Status:      "enabled",
				Type:        "default",
			}

			// negative
			client.SetToken("qqweqwe")
			res := client.NewDevice(deviceDefault)
			ctx.So(res.Code, ShouldEqual, 401)

			// login
			client.SetToken("")
			client.BasicAuth("admin@e154.ru", "admin")
			res = client.Signin()
			ctx.So(res.Code, ShouldEqual, 200)
			currentUser := &models.AuthSignInResponse{}
			err = json.Unmarshal(res.Body.Bytes(), currentUser)
			ctx.So(err, ShouldBeNil)
			accessToken = currentUser.AccessToken

			// positive
			client.SetToken(accessToken)
			res = client.NewDevice(deviceDefault)
			ctx.So(res.Code, ShouldEqual, 200)

			//TODO add node
			for i, req := range devices {
				//ctx.Println(req.Device.Name)
				client.SetToken(accessToken)
				res = client.NewDevice(req.Device)
				ctx.So(res.Code, ShouldEqual, req.ResponseCode)

				if req.ResponseCode != 200 {
					continue
				}

				device := &models.Device{}
				err = json.Unmarshal(res.Body.Bytes(), device)
				ctx.So(err, ShouldBeNil)
				devices[i].Id = device.Id
			}
		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("GET /device/{id}", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

			client := NewClient(server.GetEngine())
			client.SetToken(accessToken)

			for _, req := range devices {

				if req.Id == 0 {
					continue
				}
				//ctx.Println(req.Id)

				res := client.GetDevice(req.Id)
				ctx.So(res.Code, ShouldEqual, 200)

				if req.ResponseCode != 200 {
					continue
				}

				device := &models.Device{}
				err := json.Unmarshal(res.Body.Bytes(), device)
				ctx.So(err, ShouldBeNil)
				ctx.So(device.Id, ShouldEqual, req.Id)
				ctx.So(device.Name, ShouldEqual, req.Device.Name)
				ctx.So(device.Description, ShouldEqual, req.Device.Description)
				ctx.So(device.Type, ShouldEqual, req.Device.Type)
				ctx.So(device.Status, ShouldEqual, req.Device.Status)
				ctx.So(device.Properties, ShouldNotBeNil)
			}

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("PUT /device/{id}", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("DELETE /device/{id}", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("GET /devices", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

		})
		if err != nil {
			panic(err.Error())
		}
	})

	Convey("GET /devices/search", t, func(ctx C) {
		err := container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService *scripts.ScriptService,
			server *server.Server,
			core *core.Core) {

			err := core.Stop()
			So(err, ShouldBeNil)

			server.Shutdown()
		})
		if err != nil {
			panic(err.Error())
		}
	})
}
