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

package config

import "time"

// AppConfig ...
type AppConfig struct {
	ServerHost                     string        `json:"server_host"`
	ServerPort                     int           `json:"server_port"`
	PgUser                         string        `json:"pg_user"`
	PgPass                         string        `json:"pg_pass"`
	PgHost                         string        `json:"pg_host"`
	PgName                         string        `json:"pg_name"`
	PgPort                         string        `json:"pg_port"`
	PgDebug                        bool          `json:"pg_debug"`
	PgLogger                       bool          `json:"pg_logger"`
	PgMaxIdleConns                 int           `json:"pg_max_idle_conns"`
	PgMaxOpenConns                 int           `json:"pg_max_open_conns"`
	PgConnMaxLifeTime              int           `json:"pg_conn_max_life_time"`
	AutoMigrate                    bool          `json:"auto_migrate"`
	SnapshotDir                    string        `json:"snapshot_dir"`
	Mode                           RunMode       `json:"mode"`
	MqttPort                       int           `json:"mqtt_port"`
	MqttRetryInterval              time.Duration `json:"mqtt_retry_interval"`
	MqttRetryCheckInterval         time.Duration `json:"mqtt_retry_check_interval"`
	MqttSessionExpiryInterval      time.Duration `json:"mqtt_session_expiry_interval"`
	MqttSessionExpireCheckInterval time.Duration `json:"mqtt_session_expire_check_interval"`
	MqttQueueQos0Messages          bool          `json:"mqtt_queue_qos_0_messages"`
	MqttMaxInflight                int           `json:"mqtt_max_inflight"`
	MqttMaxAwaitRel                int           `json:"mqtt_max_await_rel"`
	MqttMaxMsgQueue                int           `json:"mqtt_max_msg_queue"`
	MqttDeliverMode                int           `json:"mqtt_deliver_mode"`
	Logging                        bool          `json:"logging"`
	Metric                         bool          `json:"metric"`
	MetricPort                     int           `json:"metric_port"`
	ColoredLogging                 bool          `json:"colored_logging"`
	AlexaHost                      string        `json:"alexa_host"`
	AlexaPort                      int           `json:"alexa_port"`
	MobileHost                     string        `json:"mobile_host"`
	MobilePort                     int           `json:"mobile_port"`
	ApiGrpcHostPort                string        `json:"api_grpc_host_port"`
	ApiHttpHostPort                string        `json:"api_http_host_port"`
	ApiPromHostPort                string        `json:"api_prom_host_port"`
	ApiWsHostPort                  string        `json:"api_ws_host_port"`
	ApiSwagger                     bool          `json:"api_swagger"`
}

// RunMode ...
type RunMode string

const (
	// DebugMode ...
	DebugMode = RunMode("debug")
	// ReleaseMode ...
	ReleaseMode = RunMode("release")
)
