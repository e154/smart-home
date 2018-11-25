package backup

import (
	"github.com/e154/smart-home/system/config"
)

type BackupConfig struct {
	Path   string
	PgUser string
	PgPass string
	PgHost string
	PgName string
	PgPort string
}

func NewBackupConfig(cfg *config.AppConfig) *BackupConfig {
	return &BackupConfig{
		Path:   cfg.SnapshotDir,
		PgUser: cfg.PgUser,
		PgPass: cfg.PgPass,
		PgHost: cfg.PgHost,
		PgName: cfg.PgName,
		PgPort: cfg.PgPort,
	}
}
