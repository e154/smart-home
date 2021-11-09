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

package models

import (
	"github.com/e154/smart-home/common"
	"time"
)

// AppConfig ...
type AppConfig struct {
	PgUser                         string         `json:"pg_user" env:"PG_USER"`
	PgPass                         string         `json:"pg_pass" env:"PG_PASS"`
	PgHost                         string         `json:"pg_host" env:"PG_HOST"`
	PgName                         string         `json:"pg_name" env:"PG_NAME"`
	PgPort                         string         `json:"pg_port" env:"PG_PORT"`
	PgDebug                        bool           `json:"pg_debug" env:"PG_DEBUG"`
	PgLogger                       bool           `json:"pg_logger" env:"PG_LOGGER"`
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
	Metric                         bool           `json:"metric" env:"METRIC"`
	MetricPort                     int            `json:"metric_port" env:"METRIC_PORT"`
	ColoredLogging                 bool           `json:"colored_logging" env:"API_GRPC_HOST_PORT"`
	AlexaHost                      string         `json:"alexa_host" env:"ALEXA_HOST"`
	AlexaPort                      int            `json:"alexa_port" env:"ALEXA_PORT"`
	MobileHost                     string         `json:"mobile_host" env:"MOBILE_HOST"`
	MobilePort                     int            `json:"mobile_port" env:"MOBILE_PORT"`
	ApiGrpcHostPort                string         `json:"api_grpc_host_port" env:"API_GRPC_HOST_PORT"`
	ApiHttpHostPort                string         `json:"api_http_host_port" env:"API_HTTP_HOST_PORT"`
	ApiPromHostPort                string         `json:"api_prom_host_port" env:"API_PROM_HOST_PORT"`
	ApiWsHostPort                  string         `json:"api_ws_host_port" env:"API_WS_HOST_PORT"`
	ApiSwagger                     bool           `json:"api_swagger" env:"API_SWAGGER"`
}
