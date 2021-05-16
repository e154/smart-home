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

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

const filenName = "config.json"

// ReadConfig ...
func ReadConfig() (conf *AppConfig, err error) {
	var file []byte
	file, err = os.ReadFile(path.Join("conf", filenName))
	if err != nil {
		log.Fatal("Error reading config file")
		return
	}
	conf = &AppConfig{}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal("Error: wrong format of config file")
		return
	}

	checkEnv(conf)

	return
}

func checkEnv(conf *AppConfig) {

	if serverHost := os.Getenv("SERVER_HOST"); serverHost != "" {
		conf.ServerHost = serverHost
	}

	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		v, _ := strconv.ParseInt(serverPort, 10, 32)
		conf.ServerPort = int(v)
	}

	if pgUser := os.Getenv("PG_USER"); pgUser != "" {
		conf.PgUser = pgUser
	}

	if pgPass := os.Getenv("PG_PASS"); pgPass != "" {
		conf.PgPass = pgPass
	}

	if pgHost := os.Getenv("PG_HOST"); pgHost != "" {
		conf.PgHost = pgHost
	}

	if pgName := os.Getenv("PG_NAME"); pgName != "" {
		conf.PgName = pgName
	}

	if pgPort := os.Getenv("PG_PORT"); pgPort != "" {
		conf.PgPort = pgPort
	}

	if pgDebug := os.Getenv("PG_DEBUG"); pgDebug != "" {
		conf.PgDebug, _ = strconv.ParseBool(pgDebug)
	}

	if pgLogger := os.Getenv("PG_LOGGER"); pgLogger != "" {
		conf.PgLogger, _ = strconv.ParseBool(pgLogger)
	}

	if pgMaxIdleConns := os.Getenv("PG_MAX_IDLE_CONNS"); pgMaxIdleConns != "" {
		v, _ := strconv.ParseInt(pgMaxIdleConns, 10, 32)
		conf.PgMaxIdleConns = int(v)
	}

	if pgMaxOpenConns := os.Getenv("PG_MAX_OPEN_CONNS"); pgMaxOpenConns != "" {
		v, _ := strconv.ParseInt(pgMaxOpenConns, 10, 32)
		conf.PgMaxOpenConns = int(v)
	}

	if pgConnMaxLifeTime := os.Getenv("PG_CONN_MAX_LIFE_TIME"); pgConnMaxLifeTime != "" {
		v, _ := strconv.ParseInt(pgConnMaxLifeTime, 10, 32)
		conf.PgConnMaxLifeTime = int(v)
	}

	if autoMigrate := os.Getenv("AUTO_MIGRATE"); autoMigrate != "" {
		conf.AutoMigrate, _ = strconv.ParseBool(autoMigrate)
	}

	if snapshotDir := os.Getenv("SNAPSHOT_DIR"); snapshotDir != "" {
		conf.SnapshotDir = snapshotDir
	}

	if mode := os.Getenv("MODE"); mode != "" {
		conf.Mode = RunMode(mode)
	}

	if mqttPort := os.Getenv("MQTT_PORT"); mqttPort != "" {
		v, _ := strconv.ParseInt(mqttPort, 10, 32)
		conf.MqttPort = int(v)
	}

	if mqttRetryInterval := os.Getenv("MQTT_RETRY_INTERVAL"); mqttRetryInterval != "" {
		v, _ := strconv.ParseInt(mqttRetryInterval, 10, 32)
		conf.MqttRetryInterval = time.Duration(v)
	}

	if mqttRetryCheckInterval := os.Getenv("MQTT_RETRY_CHECK_INTERVAL"); mqttRetryCheckInterval != "" {
		v, _ := strconv.ParseInt(mqttRetryCheckInterval, 10, 32)
		conf.MqttRetryCheckInterval = time.Duration(v)
	}

	if mqttSessionExpiryInterval := os.Getenv("MQTT_SESSION_EXPIRY_INTERVAL"); mqttSessionExpiryInterval != "" {
		v, _ := strconv.ParseInt(mqttSessionExpiryInterval, 10, 32)
		conf.MqttSessionExpiryInterval = time.Duration(v)
	}

	if mqttSessionExpireCheckInterval := os.Getenv("MQTT_SESSION_EXPIRE_CHECK_INTERVAL"); mqttSessionExpireCheckInterval != "" {
		v, _ := strconv.ParseInt(mqttSessionExpireCheckInterval, 10, 32)
		conf.MqttSessionExpireCheckInterval = time.Duration(v)
	}

	if mqttQueueQos0Messages := os.Getenv("MQTT_QUEUE_QOS_0_MESSAGES"); mqttQueueQos0Messages != "" {
		conf.MqttQueueQos0Messages, _ = strconv.ParseBool(mqttQueueQos0Messages)
	}

	if mqttKeepAlive := os.Getenv("MQTT_MAX_INFLIGHT"); mqttKeepAlive != "" {
		v, _ := strconv.ParseInt(mqttKeepAlive, 10, 32)
		conf.MqttMaxInflight = int(v)
	}

	if mqttMaxAwaitRel := os.Getenv("MQTT_MAX_AWAIT_REL"); mqttMaxAwaitRel != "" {
		v, _ := strconv.ParseInt(mqttMaxAwaitRel, 10, 32)
		conf.MqttMaxAwaitRel = int(v)
	}

	if mqttMaxMsgQueue := os.Getenv("MQTT_MAX_MSG_QUEUE"); mqttMaxMsgQueue != "" {
		v, _ := strconv.ParseInt(mqttMaxMsgQueue, 10, 32)
		conf.MqttMaxMsgQueue = int(v)
	}

	if mqttDeliverMode := os.Getenv("MQTT_DELIVER_MODE"); mqttDeliverMode != "" {
		v, _ := strconv.ParseInt(mqttDeliverMode, 10, 32)
		conf.MqttDeliverMode = int(v)
	}

	if logging := os.Getenv("LOGGING"); logging != "" {
		conf.Logging, _ = strconv.ParseBool(logging)
	}

	if metricPort := os.Getenv("METRIC_PORT"); metricPort != "" {
		v, _ := strconv.ParseInt(metricPort, 10, 32)
		conf.MetricPort = int(v)
	}
}
