// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"os"
	"testing"

	"github.com/e154/smart-home/system/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {

	t.Run("config", func(t *testing.T) {
		Convey("", t, func(ctx C) {

			err := container.Invoke(func() {

				t.Run("file", func(t *testing.T) {
					Convey("", t, func(ctx C) {
						conf, err := config.ReadConfig("tests/data", "config.json", "")
						ctx.So(err, ShouldBeNil)

						//debug.Println(conf)

						ctx.So(conf.PgUser, ShouldEqual, "smart_home")
						ctx.So(conf.PgPass, ShouldEqual, "smart_home")
						ctx.So(conf.PgHost, ShouldEqual, "127.0.0.1")
						ctx.So(conf.PgName, ShouldEqual, "smart_home_test")
						ctx.So(conf.PgPort, ShouldEqual, "5432")
						ctx.So(conf.PgDebug, ShouldEqual, true)
						ctx.So(conf.PgLogger, ShouldEqual, true)
						ctx.So(conf.PgMaxIdleConns, ShouldEqual, 10)
						ctx.So(conf.PgMaxOpenConns, ShouldEqual, 50)
						ctx.So(conf.PgConnMaxLifeTime, ShouldEqual, 30)
						ctx.So(conf.AutoMigrate, ShouldEqual, true)
						ctx.So(conf.SnapshotDir, ShouldEqual, "snapshots")
						ctx.So(conf.Mode, ShouldEqual, "release")
						ctx.So(conf.MqttPort, ShouldEqual, 1883)
						ctx.So(conf.MqttRetryInterval, ShouldEqual, 20)
						ctx.So(conf.MqttRetryCheckInterval, ShouldEqual, 20)
						ctx.So(conf.MqttSessionExpiryInterval, ShouldEqual, 0)
						ctx.So(conf.MqttSessionExpireCheckInterval, ShouldEqual, 0)
						ctx.So(conf.MqttQueueQos0Messages, ShouldEqual, false)
						ctx.So(conf.MqttMaxInflight, ShouldEqual, 32)
						ctx.So(conf.MqttMaxAwaitRel, ShouldEqual, 100)
						ctx.So(conf.MqttMaxMsgQueue, ShouldEqual, 1000)
						ctx.So(conf.MqttDeliverMode, ShouldEqual, 1)
						ctx.So(conf.Logging, ShouldEqual, true)
						ctx.So(conf.ColoredLogging, ShouldEqual, false)
						ctx.So(conf.AlexaHost, ShouldEqual, "0.0.0.0")
						ctx.So(conf.AlexaPort, ShouldEqual, 3002)
						ctx.So(conf.MobileHost, ShouldEqual, "0.0.0.0")
						ctx.So(conf.MobilePort, ShouldEqual, 3001)
						ctx.So(conf.ApiGrpcPort, ShouldEqual, 3000)
						ctx.So(conf.ApiHttpPort, ShouldEqual, 3001)
						ctx.So(conf.ApiSwagger, ShouldEqual, true)
						ctx.So(conf.Domain, ShouldEqual, "localhost")
						ctx.So(conf.Https, ShouldEqual, false)
					})
				})

				t.Run("env", func(t *testing.T) {
					Convey("", t, func(ctx C) {

						_ = os.Setenv("PG_USER", "_smart_home")
						_ = os.Setenv("PG_USER2", "smart_home")
						_ = os.Setenv("PG_PASS", "_smart_home")
						_ = os.Setenv("PG_HOST", "127.0.0.1")
						_ = os.Setenv("PG_NAME", "smart_home")
						_ = os.Setenv("PG_PORT", "5432")
						_ = os.Setenv("PG_DEBUG", "false")
						_ = os.Setenv("PG_LOGGER", "false")
						_ = os.Setenv("PG_MAX_IDLE_CONNS", "10")
						_ = os.Setenv("PG_MAX_OPEN_CONNS", "50")
						_ = os.Setenv("PG_CONN_MAX_LIFE_TIME", "30")
						_ = os.Setenv("AUTO_MIGRATE", "true")
						_ = os.Setenv("SNAPSHOT_DIR", "snapshots")
						_ = os.Setenv("MODE", "release")
						_ = os.Setenv("MQTT_PORT", "1883")
						_ = os.Setenv("MQTT_RETRY_INTERVAL", "20")
						_ = os.Setenv("MQTT_RETRY_CHECK_INTERVAL", "20")
						_ = os.Setenv("MQTT_SESSION_EXPIRY_INTERVAL", "0")
						_ = os.Setenv("MQTT_SESSION_EXPIRE_CHECK_INTERVAL", "0")
						_ = os.Setenv("MQTT_QUEUE_QOS_0_MESSAGES", "true")
						_ = os.Setenv("MQTT_MAX_INFLIGHT", "32")
						_ = os.Setenv("MQTT_MAX_AWAIT_REL", "100")
						_ = os.Setenv("MQTT_MAX_MSG_QUEUE", "1000")
						_ = os.Setenv("MQTT_DELIVER_MODE", "1")
						_ = os.Setenv("LOGGING", "true")
						_ = os.Setenv("METRIC", "false")
						_ = os.Setenv("METRIC_PORT", "2112")
						_ = os.Setenv("API_GRPC_HOST_PORT", ":3000")
						_ = os.Setenv("ALEXA_HOST", "0.0.0.0")
						_ = os.Setenv("ALEXA_PORT", "3002")
						_ = os.Setenv("MOBILE_HOST", "0.0.0.0")
						_ = os.Setenv("MOBILE_PORT", "3001")
						_ = os.Setenv("API_GRPC_HOST_PORT", ":3000")
						_ = os.Setenv("API_HTTP_HOST_PORT", ":3001")
						_ = os.Setenv("API_WS_HOST_PORT", ":3003")
						_ = os.Setenv("API_SWAGGER", "true")
						_ = os.Setenv("DOMAIN", "localhost")
						_ = os.Setenv("HTTPS", "false")

						conf, err := config.ReadConfig("tests/data", "config.json", "")
						ctx.So(err, ShouldBeNil)

						//debug.Println(conf)

						ctx.So(conf.PgUser, ShouldEqual, "_smart_home")
						ctx.So(conf.PgPass, ShouldEqual, "_smart_home")
						ctx.So(conf.PgHost, ShouldEqual, "127.0.0.1")
						ctx.So(conf.PgName, ShouldEqual, "smart_home")
						ctx.So(conf.PgPort, ShouldEqual, "5432")
						ctx.So(conf.PgDebug, ShouldEqual, false)
						ctx.So(conf.PgLogger, ShouldEqual, false)
						ctx.So(conf.PgMaxIdleConns, ShouldEqual, 10)
						ctx.So(conf.PgMaxOpenConns, ShouldEqual, 50)
						ctx.So(conf.PgConnMaxLifeTime, ShouldEqual, 30)
						ctx.So(conf.AutoMigrate, ShouldEqual, true)
						ctx.So(conf.SnapshotDir, ShouldEqual, "snapshots")
						ctx.So(conf.Mode, ShouldEqual, "release")
						ctx.So(conf.MqttPort, ShouldEqual, 1883)
						ctx.So(conf.MqttRetryInterval, ShouldEqual, 20)
						ctx.So(conf.MqttRetryCheckInterval, ShouldEqual, 20)
						ctx.So(conf.MqttSessionExpiryInterval, ShouldEqual, 0)
						ctx.So(conf.MqttSessionExpireCheckInterval, ShouldEqual, 0)
						ctx.So(conf.MqttQueueQos0Messages, ShouldEqual, true)
						ctx.So(conf.MqttMaxInflight, ShouldEqual, 32)
						ctx.So(conf.MqttMaxAwaitRel, ShouldEqual, 100)
						ctx.So(conf.MqttMaxMsgQueue, ShouldEqual, 1000)
						ctx.So(conf.MqttDeliverMode, ShouldEqual, 1)
						ctx.So(conf.Logging, ShouldEqual, true)
						ctx.So(conf.ColoredLogging, ShouldEqual, false)
						ctx.So(conf.AlexaHost, ShouldEqual, "0.0.0.0")
						ctx.So(conf.AlexaPort, ShouldEqual, 3002)
						ctx.So(conf.MobileHost, ShouldEqual, "0.0.0.0")
						ctx.So(conf.MobilePort, ShouldEqual, 3001)
						ctx.So(conf.ApiGrpcPort, ShouldEqual, 3000)
						ctx.So(conf.ApiHttpPort, ShouldEqual, 3001)
						ctx.So(conf.ApiSwagger, ShouldEqual, true)
						ctx.So(conf.Domain, ShouldEqual, "localhost")
						ctx.So(conf.Https, ShouldEqual, false)
					})
				})

				t.Run("env + prefix", func(t *testing.T) {
					Convey("", t, func(ctx C) {

						_ = os.Setenv("APP_PG_USER", "_smart_home_")
						_ = os.Setenv("APP_PG_USER2", "smart_home2")
						_ = os.Setenv("APP_PG_PASS", "_smart_home_")
						_ = os.Setenv("APP_PG_HOST", "127.0.0.1")
						_ = os.Setenv("APP_PG_NAME", "smart_home")
						_ = os.Setenv("APP_PG_PORT", "5432")
						_ = os.Setenv("APP_PG_DEBUG", "false")
						_ = os.Setenv("APP_PG_LOGGER", "false")
						_ = os.Setenv("APP_PG_MAX_IDLE_CONNS", "10")
						_ = os.Setenv("APP_PG_MAX_OPEN_CONNS", "50")
						_ = os.Setenv("APP_PG_CONN_MAX_LIFE_TIME", "30")
						_ = os.Setenv("APP_AUTO_MIGRATE", "true")
						_ = os.Setenv("APP_SNAPSHOT_DIR", "snapshots")
						_ = os.Setenv("APP_MODE", "release")
						_ = os.Setenv("APP_MQTT_PORT", "1883")
						_ = os.Setenv("APP_MQTT_RETRY_INTERVAL", "20")
						_ = os.Setenv("APP_MQTT_RETRY_CHECK_INTERVAL", "20")
						_ = os.Setenv("APP_MQTT_SESSION_EXPIRY_INTERVAL", "0")
						_ = os.Setenv("APP_MQTT_SESSION_EXPIRE_CHECK_INTERVAL", "0")
						_ = os.Setenv("APP_MQTT_QUEUE_QOS_0_MESSAGES", "true")
						_ = os.Setenv("APP_MQTT_MAX_INFLIGHT", "32")
						_ = os.Setenv("APP_MQTT_MAX_AWAIT_REL", "100")
						_ = os.Setenv("APP_MQTT_MAX_MSG_QUEUE", "1000")
						_ = os.Setenv("APP_MQTT_DELIVER_MODE", "1")
						_ = os.Setenv("APP_LOGGING", "true")
						_ = os.Setenv("APP_METRIC", "false")
						_ = os.Setenv("APP_METRIC_PORT", "2112")
						_ = os.Setenv("APP_API_GRPC_HOST_PORT", ":3000")
						_ = os.Setenv("APP_ALEXA_HOST", "0.0.0.0")
						_ = os.Setenv("APP_ALEXA_PORT", "3002")
						_ = os.Setenv("APP_MOBILE_HOST", "0.0.0.0")
						_ = os.Setenv("APP_MOBILE_PORT", "3001")
						_ = os.Setenv("APP_API_GRPC_HOST_PORT", ":3000")
						_ = os.Setenv("APP_API_HTTP_HOST_PORT", ":3001")
						_ = os.Setenv("APP_API_WS_HOST_PORT", ":3003")
						_ = os.Setenv("APP_API_SWAGGER", "true")
						_ = os.Setenv("DOMAIN", "localhost")
						_ = os.Setenv("HTTPS", "false")

						conf, err := config.ReadConfig("tests/data", "config.json", "APP")
						ctx.So(err, ShouldBeNil)

						//debug.Println(conf)

						ctx.So(conf.PgUser, ShouldEqual, "_smart_home_")
						ctx.So(conf.PgPass, ShouldEqual, "_smart_home_")
						ctx.So(conf.PgHost, ShouldEqual, "127.0.0.1")
						ctx.So(conf.PgName, ShouldEqual, "smart_home")
						ctx.So(conf.PgPort, ShouldEqual, "5432")
						ctx.So(conf.PgDebug, ShouldEqual, false)
						ctx.So(conf.PgLogger, ShouldEqual, false)
						ctx.So(conf.PgMaxIdleConns, ShouldEqual, 10)
						ctx.So(conf.PgMaxOpenConns, ShouldEqual, 50)
						ctx.So(conf.PgConnMaxLifeTime, ShouldEqual, 30)
						ctx.So(conf.AutoMigrate, ShouldEqual, true)
						ctx.So(conf.SnapshotDir, ShouldEqual, "snapshots")
						ctx.So(conf.Mode, ShouldEqual, "release")
						ctx.So(conf.MqttPort, ShouldEqual, 1883)
						ctx.So(conf.MqttRetryInterval, ShouldEqual, 20)
						ctx.So(conf.MqttRetryCheckInterval, ShouldEqual, 20)
						ctx.So(conf.MqttSessionExpiryInterval, ShouldEqual, 0)
						ctx.So(conf.MqttSessionExpireCheckInterval, ShouldEqual, 0)
						ctx.So(conf.MqttQueueQos0Messages, ShouldEqual, true)
						ctx.So(conf.MqttMaxInflight, ShouldEqual, 32)
						ctx.So(conf.MqttMaxAwaitRel, ShouldEqual, 100)
						ctx.So(conf.MqttMaxMsgQueue, ShouldEqual, 1000)
						ctx.So(conf.MqttDeliverMode, ShouldEqual, 1)
						ctx.So(conf.Logging, ShouldEqual, true)
						ctx.So(conf.ColoredLogging, ShouldEqual, false)
						ctx.So(conf.AlexaHost, ShouldEqual, "0.0.0.0")
						ctx.So(conf.AlexaPort, ShouldEqual, 3002)
						ctx.So(conf.MobileHost, ShouldEqual, "0.0.0.0")
						ctx.So(conf.MobilePort, ShouldEqual, 3001)
						ctx.So(conf.ApiGrpcPort, ShouldEqual, 3000)
						ctx.So(conf.ApiHttpPort, ShouldEqual, 3001)
						ctx.So(conf.ApiSwagger, ShouldEqual, true)
						ctx.So(conf.Domain, ShouldEqual, "localhost")
						ctx.So(conf.Https, ShouldEqual, false)
					})
				})
			})
			ctx.So(err, ShouldBeNil)
		})
	})
}
