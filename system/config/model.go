package config

import "time"

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
}

type RunMode string

const (
	DebugMode   = RunMode("debug")
	ReleaseMode = RunMode("release")
)
