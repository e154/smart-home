package config

type AppConfig struct {
	ServerHost           string  `json:"server_host"`
	ServerPort           int     `json:"server_port"`
	PgUser               string  `json:"pg_user"`
	PgPass               string  `json:"pg_pass"`
	PgHost               string  `json:"pg_host"`
	PgName               string  `json:"pg_name"`
	PgPort               string  `json:"pg_port"`
	PgDebug              bool    `json:"pg_debug"`
	PgLogger             bool    `json:"pg_logger"`
	PgMaxIdleConns       int     `json:"pg_max_idle_conns"`
	PgMaxOpenConns       int     `json:"pg_max_open_conns"`
	PgConnMaxLifeTime    int     `json:"pg_conn_max_life_time"`
	AutoMigrate          bool    `json:"auto_migrate"`
	SnapshotDir          string  `json:"snapshot_dir"`
	Mode                 RunMode `json:"mode"`
	MqttKeepAlive        int     `json:"mqtt_keep_alive"`
	MqttConnectTimeout   int     `json:"mqtt_connect_timeout"`
	MqttSessionsProvider string  `json:"mqtt_sessions_provider"`
	MqttTopicsProvider   string  `json:"mqtt_topics_provider"`
	MqttPort             int     `json:"mqtt_port"`
}

type RunMode string

const (
	DebugMode   = RunMode("debug")
	ReleaseMode = RunMode("release")
)
