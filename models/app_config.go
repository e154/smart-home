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

package models

import (
	"fmt"
	"time"

	"github.com/e154/smart-home/common"
)

// AppConfig ...
type AppConfig struct {
	PgUser                         string         `json:"pg_user" env:"PG_USER"`
	PgPass                         string         `json:"pg_pass" env:"PG_PASS"`
	PgHost                         string         `json:"pg_host" env:"PG_HOST"`
	PgName                         string         `json:"pg_name" env:"PG_NAME"`
	PgPort                         string         `json:"pg_port" env:"PG_PORT"`
	PgDebug                        bool           `json:"pg_debug" env:"PG_DEBUG"`
	PgMaxIdleConns                 int            `json:"pg_max_idle_conns" env:"PG_MAX_IDLE_CONNS"`
	PgMaxOpenConns                 int            `json:"pg_max_open_conns" env:"PG_MAX_OPEN_CONNS"`
	PgConnMaxLifeTime              int            `json:"pg_conn_max_life_time" env:"PG_CONN_MAX_LIFE_TIME"`
	AutoMigrate                    bool           `json:"auto_migrate" env:"AUTO_MIGRATE"`
	SnapshotDir                    string         `json:"snapshot_dir" env:"SNAPSHOT_DIR"`
	Mode                           common.RunMode `json:"mode" env:"MODE"`
	MqttPort                       int            `json:"mqtt_port" env:"MQTT_PORT"`
	MqttRetryInterval              time.Duration  `json:"mqtt_retry_interval" env:"MQTT_RETRY_INTERVAL"`
	MqttRetryCheckInterval         time.Duration  `json:"mqtt_retry_check_interval" env:"MQTT_RETRY_CHECK_INTERVAL"`
	MqttSessionExpiryInterval      time.Duration  `json:"mqtt_session_expiry_interval" env:"MQTT_SESSION_EXPIRY_INTERVAL"`
	MqttSessionExpireCheckInterval time.Duration  `json:"mqtt_session_expire_check_interval" env:"MQTT_SESSION_EXPIRE_CHECK_INTERVAL"`
	MqttQueueQos0Messages          bool           `json:"mqtt_queue_qos_0_messages" env:"MQTT_QUEUE_QOS_0_MESSAGES"`
	MqttMaxInflight                int            `json:"mqtt_max_inflight" env:"MQTT_MAX_INFLIGHT"`
	MqttMaxAwaitRel                int            `json:"mqtt_max_await_rel" env:"MQTT_MAX_AWAIT_REL"`
	MqttMaxMsgQueue                int            `json:"mqtt_max_msg_queue" env:"MQTT_MAX_MSG_QUEUE"`
	MqttDeliverMode                int            `json:"mqtt_deliver_mode" env:"MQTT_DELIVER_MODE"`
	Logging                        bool           `json:"logging" env:"LOGGING"`
	ColoredLogging                 bool           `json:"colored_logging" env:"COLORED_LOGGING"`
	AlexaHost                      string         `json:"alexa_host" env:"ALEXA_HOST"`
	AlexaPort                      int            `json:"alexa_port" env:"ALEXA_PORT"`
	MobileHost                     string         `json:"mobile_host" env:"MOBILE_HOST"`
	MobilePort                     int            `json:"mobile_port" env:"MOBILE_PORT"`
	ApiGrpcPort                    int            `json:"api_grpc_port" env:"API_GRPC_PORT"`
	ApiHttpPort                    int            `json:"api_http_port" env:"API_HTTP_PORT"`
	ApiSwagger                     bool           `json:"api_swagger" env:"API_SWAGGER"`
	Lang                           string         `json:"lang" env:"LANG"`
	GodMode                        bool           `json:"god_mode" env:"GOD_MODE"`
	Domain                         string         `json:"domain" env:"DOMAIN"`
	Https                          bool           `json:"https" env:"HTTPS"`
}

func (c *AppConfig) ApiScheme() (scheme string) {
	scheme = "http"
	if c.Https {
		scheme = "https"
	}
	return
}

func (c *AppConfig) ApiFullAddress() (scheme string) {
	return fmt.Sprintf("%s://%s:%d", c.ApiScheme(), c.Domain, c.ApiHttpPort)
}
