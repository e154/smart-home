package api

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/server"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/api/server/v1/responses"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestDevice(t *testing.T) {

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

	type updateDeviceRequest struct {
		ResponseCode int
		Device       models.UpdateDevice
	}

	updateDevices := []updateDeviceRequest{
		{
			ResponseCode: 200,
			Device: models.UpdateDevice{
				Name: "device1", Description: "device desc", Status: "enabled", Type: "default",
			},
		},
		{
			ResponseCode: 200,
			Device: models.UpdateDevice{
				Name: "device2", Description: "device2 desc", Status: "disabled", Type: "modbus_rtu",
				Properties: models.DeviceProperties{
					DevModBusRtuConfig: &models.DevModBusRtuConfig{
						SlaveId: 1, Baud: 115200, DataBits: 8, StopBits: 2, Parity: "none", Timeout: 100,
					},
				},
			},
		},
		{
			ResponseCode: 200,
			Device: models.UpdateDevice{
				Name: "device1", Description: "device2 desc", Status: "enabled", Type: "modbus_tcp",
				Properties: models.DeviceProperties{
					DevModBusTcpConfig: &models.DevModBusTcpConfig{
						SlaveId: 1, AddressPort: "127.0.0.1:502",
					},
				},
			},
		},
		{
			ResponseCode: 500,
			Device: models.UpdateDevice{
				Name: "device1", Description: "device2 desc", Status: "custom", Type: "modbus_tcp",
				Properties: models.DeviceProperties{
					DevModBusTcpConfig: &models.DevModBusTcpConfig{
						SlaveId: 1, AddressPort: "127.0.0.1:502",
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

			// negative
			client.SetToken(invalidToken1)
			res := client.GetDevice(1)
			ctx.So(res.Code, ShouldEqual, 401)
			client.SetToken(invalidToken2)
			res = client.GetDevice(1)
			ctx.So(res.Code, ShouldEqual, 403)

			// negative
			client.SetToken(accessToken)
			res = client.GetDevice(404)
			ctx.So(res.Code, ShouldEqual, 404)

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

			client := NewClient(server.GetEngine())

			device := &models.UpdateDevice{
				Name: "device1", Description: "device desc", Status: "enabled", Type: "modbus_rtu",
				Properties: models.DeviceProperties{
					DevModBusRtuConfig: &models.DevModBusRtuConfig{
						SlaveId: 1, Baud: 115200, DataBits: 8, StopBits: 2, Parity: "none", Timeout: 100,
					},
				},
			}

			// negative
			client.SetToken(invalidToken1)
			res := client.UpdateDevice(1, device)
			ctx.So(res.Code, ShouldEqual, 401)
			client.SetToken(invalidToken2)
			res = client.UpdateDevice(1, device)
			ctx.So(res.Code, ShouldEqual, 403)

			//
			client.SetToken(accessToken)
			res = client.UpdateDevice(404, device)
			ctx.So(res.Code, ShouldEqual, 404)

			for _, device := range updateDevices {
				res := client.UpdateDevice(1, device.Device)
				ctx.So(res.Code, ShouldEqual, device.ResponseCode)

				res = client.GetDevice(1)
				ctx.So(res.Code, ShouldEqual, 200)

				d := &models.UpdateDevice{}
				err := json.Unmarshal(res.Body.Bytes(), d)
				ctx.So(err, ShouldBeNil)
				ctx.So(d.Name, ShouldEqual, device.Device.Name)
				ctx.So(d.Description, ShouldEqual, device.Device.Description)
				ctx.So(d.Type, ShouldEqual, device.Device.Type)
				ctx.So(d.Status, ShouldBeIn, []string{device.Device.Status, "enabled"})
			}
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

			client := NewClient(server.GetEngine())

			// negative
			client.SetToken(invalidToken1)
			res := client.DeleteDevice(404)
			ctx.So(res.Code, ShouldEqual, 401)
			client.SetToken(invalidToken2)
			res = client.DeleteDevice(404)
			ctx.So(res.Code, ShouldEqual, 403)

			//
			client.SetToken(accessToken)
			res = client.DeleteDevice(404)
			ctx.So(res.Code, ShouldEqual, 404)

			res = client.DeleteDevice(1)
			ctx.So(res.Code, ShouldEqual, 200)

			res = client.DeleteDevice(1)
			ctx.So(res.Code, ShouldEqual, 404)
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

			client := NewClient(server.GetEngine())

			// negative
			client.SetToken(invalidToken1)
			res := client.GetDeviceList(5, 0, "DESC", "id")
			ctx.So(res.Code, ShouldEqual, 401)
			client.SetToken(invalidToken2)
			res = client.GetDeviceList(5, 0, "DESC", "id")
			ctx.So(res.Code, ShouldEqual, 403)

			// positive
			client.SetToken(accessToken)

			listGetter := func(limit, offset, realLimit, realOffset int) {
				res = client.GetDeviceList(limit, offset, "DESC", "id")
				ctx.So(res.Code, ShouldEqual, 200)

				deviceList := responses.DeviceList{}
				err := json.Unmarshal(res.Body.Bytes(), &deviceList.Body)
				ctx.So(err, ShouldBeNil)

				ctx.So(len(deviceList.Body.Items), ShouldEqual, realLimit)
				ctx.So(deviceList.Body.Meta.Limit, ShouldEqual, limit)
				ctx.So(deviceList.Body.Meta.Offset, ShouldEqual, realOffset)
			}

			listGetter(5, 0, 5, 0)
			listGetter(1, 3, 1, 3)
			listGetter(7, 0, 6, 0)

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
