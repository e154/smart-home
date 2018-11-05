package config

type AppConfig struct {
	ServerHost        string `json:"server_host"`
	ServerPort        int    `json:"server_port"`
	PgUser            string `json:"pg_user"`
	PgPass            string `json:"pg_pass"`
	PgHost            string `json:"pg_host"`
	PgName            string `json:"pg_name"`
	PgPort            string `json:"pg_port"`
	PgDebug           bool   `json:"pg_debug"`
	PgLogger          bool   `json:"pg_logger"`
	PgMaxIdleConns    int    `json:"pg_max_idle_conns"`
	PgMaxOpenConns    int    `json:"pg_max_open_conns"`
	PgConnMaxLifeTime int    `json:"pg_conn_max_life_time"`
	AutoMigrate       bool   `json:"auto_migrate"`
}
