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

package backup

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/fx"
	"path/filepath"

	app "github.com/e154/smart-home/common/app"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
)

var (
	log = logger.MustGetLogger("backup")
)

// Backup ...
type Backup struct {
	cfg          *Config
	restoreImage string
}

// NewBackup ...
func NewBackup(lc fx.Lifecycle, cfg *Config) *Backup {

	backup := &Backup{
		cfg: cfg,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return backup.Shutdown(ctx)
		},
	})

	return backup
}

// Shutdown ...
func (b *Backup) Shutdown(ctx context.Context) (err error) {

	if b.restoreImage != "" {
		if err = b.restore(b.restoreImage); err != nil {
			log.Errorf("%+v", err)
			return
		}
		app.IsRestart = true
	}
	return
}

// New ...
func (b *Backup) New() (err error) {
	log.Info("create new backup")

	tmpDir := path.Join(os.TempDir(), "smart_home")
	if err = os.MkdirAll(tmpDir, 0755); err != nil {
		return
	}

	if err = NewLocal(b.cfg).New(tmpDir); err != nil {
		log.Error(err.Error())
		return
	}

	err = zipit([]string{
		path.Join("data", "file_storage"),
		path.Join(tmpDir, "data.sql"),
		path.Join(tmpDir, "scheme.sql"),
	},
		path.Join(b.cfg.Path, fmt.Sprintf("%s.zip", time.Now().Format("2006-01-02T15:04:05.999"))))
	if err != nil {
		return
	}

	_ = os.RemoveAll(tmpDir)

	log.Info("complete")

	return
}

// List ...
func (b *Backup) List() (list []string) {

	_ = filepath.Walk(b.cfg.Path, func(path string, info os.FileInfo, err error) error {
		if info.Name() == ".gitignore" || info.Name() == b.cfg.Path || info.IsDir() {
			return nil
		}
		if info.Name()[0:1] == "." {
			return nil
		}
		list = append(list, info.Name())
		return nil
	})
	return
}

// Restore ...
func (b *Backup) Restore(name string) (err error) {
	if name == "" {
		return
	}

	b.restoreImage = name

	log.Info("try to shutdown")
	err = app.Kill()

	return
}

// restore ...
func (b *Backup) restore(name string) (err error) {
	log.Infof("restore backup file %s", name)

	file := path.Join(b.cfg.Path, name)

	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("path %s", file))
		return
	}

	tmpDir := path.Join(os.TempDir(), "smart_home")
	if err = unzip(file, tmpDir); err != nil {
		err = errors.Wrap(fmt.Errorf("failed unzip file %s", file), err.Error())
		return
	}

	log.Info("drop database")

	var db *gorm.DB
	db, err = gorm.Open(postgres.Open(b.cfg.String()), &gorm.Config{})
	if err != nil {
		return
	}
	defer func() {
		if db, err := db.DB(); err == nil {
			_ = db.Close()
		}
	}()

	if err = db.Exec(`DROP SCHEMA IF EXISTS "public" CASCADE;`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed exec sql command"), err.Error())
		return
	}

	if err = db.Exec(`CREATE SCHEMA "public";`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed exec sql command"), err.Error())
		return
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;`)

	log.Info("restore database scheme")
	var filename = path.Join(tmpDir, "scheme.sql")
	if err = NewLocal(b.cfg).Restore(filename); err != nil {
		return
	}

	db.Exec(`SELECT create_hypertable('metric_bucket', 'time', migrate_data => true, if_not_exists => TRUE);`)

	log.Info("restore database data")
	filename = path.Join(tmpDir, "data.sql")
	if err = NewLocal(b.cfg).Restore(filename); err != nil {
		return
	}

	log.Info("restore files ...")
	d := path.Join("data", "file_storage")
	log.Infof("remove data dir")
	_ = os.RemoveAll(d)

	from := path.Join(tmpDir, "file_storage")
	to := path.Join("data", "file_storage")
	log.Infof("copy file_storage %s --> %s", from, to)
	if err = Copy(from, to); err != nil {
		return
	}

	log.Infof("remove tmp dir %s", tmpDir)
	_ = os.RemoveAll(tmpDir)

	log.Info("complete ...")

	return
}
