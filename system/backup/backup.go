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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
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
	eventBus     bus.Bus
	restoreImage string
}

// NewBackup ...
func NewBackup(lc fx.Lifecycle, eventBus bus.Bus, cfg *Config) *Backup {

	if cfg.Path == "" {
		cfg.Path = "snapshots"
	}

	backup := &Backup{
		cfg:      cfg,
		eventBus: eventBus,
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

	backupName := fmt.Sprintf("%s.zip", time.Now().Format("2006-01-02T15:04:05.999"))
	err = zipit([]string{
		path.Join("data", "file_storage"),
		path.Join(tmpDir, "data.sql"),
		path.Join(tmpDir, "scheme.sql"),
	},
		path.Join(b.cfg.Path, backupName))
	if err != nil {
		return
	}

	_ = os.RemoveAll(tmpDir)

	log.Info("complete")

	b.eventBus.Publish("system/services/backup", events.EventCreatedBackup{
		Name: backupName,
	})

	return
}

// List ...
func (b *Backup) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.Backup, total int64, err error) {

	_ = filepath.Walk(b.cfg.Path, func(path string, info os.FileInfo, err error) error {
		if info.Name() == ".gitignore" || info.Name() == b.cfg.Path || info.IsDir() {
			return nil
		}
		if info.Name()[0:1] == "." {
			return nil
		}
		list = append(list, &m.Backup{
			Name:     info.Name(),
			Size:     info.Size(),
			FileMode: info.Mode(),
			ModTime:  info.ModTime(),
		})
		return nil
	})
	//todo: add pagination
	return
}

// Restore ...
func (b *Backup) Restore(name string) (err error) {
	if name == "" {
		return
	}

	b.restoreImage = name
	app.Restore = true

	log.Info("try to shutdown")
	err = app.Kill()

	return
}

func (b *Backup) RollbackChanges() (err error) {

	var db *gorm.DB
	db, err = gorm.Open(postgres.Open(b.cfg.String()), &gorm.Config{})
	if err != nil {
		return
	}

	defer func() {
		if _db, err := db.DB(); err == nil {
			_ = _db.Close()
		}
	}()

	if err = db.Exec(`DROP EXTENSION IF EXISTS timescaledb CASCADE;`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed drop extension timescaledb"), err.Error())
		return
	}

	if err = db.Exec(`ALTER SCHEMA "public_old" RENAME TO "public";`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed rename public scheme"), err.Error())
		return
	}

	if err = db.Exec(`DROP SCHEMA IF EXISTS "public_old" CASCADE;`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed drop schema"), err.Error())
		return
	}

	if err = db.Exec(`CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed create extension if not exists timescaledb"), err.Error())
	}

	return
}

func (b *Backup) ApplyChanges() (err error) {

	var db *gorm.DB
	db, err = gorm.Open(postgres.Open(b.cfg.String()), &gorm.Config{})
	if err != nil {
		return
	}

	defer func() {
		if _db, err := db.DB(); err == nil {
			_ = _db.Close()
		}
	}()

	if err = db.Exec(`DROP EXTENSION IF EXISTS timescaledb CASCADE;`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed drop extension timescaledb"), err.Error())
		return
	}

	if err = db.Exec(`DROP SCHEMA IF EXISTS "public_old" CASCADE;`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed drop schema"), err.Error())
		return
	}

	if err = db.Exec(`CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed create extension if not exists timescaledb"), err.Error())
	}

	return
}

// restore ...
func (b *Backup) restore(name string) (err error) {
	log.Infof("restore backup file %s", name)

	file := path.Join(b.cfg.Path, name)

	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		err = errors.Wrap(apperr.ErrBackupNotFound, fmt.Sprintf("path %s", file))
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

		if _db, err := db.DB(); err == nil {
			_ = _db.Close()
		}
	}()

	if err = db.Exec(`DROP EXTENSION IF EXISTS timescaledb CASCADE;`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed drop extension timescaledb"), err.Error())
		return
	}

	if err = db.Exec(`CREATE SCHEMA IF NOT EXISTS "public";`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed create public schema"), err.Error())
		return
	}

	if err = db.Exec(`DROP SCHEMA IF EXISTS "public_old" CASCADE;`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed drop public_old schema"), err.Error())
		return
	}

	if err = db.Exec(`ALTER SCHEMA "public" RENAME TO "public_old";`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed rename public scheme"), err.Error())
		return
	}

	if err = db.Exec(`CREATE SCHEMA IF NOT EXISTS "public";`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed create public schema"), err.Error())
		return
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;`)
	db.Exec(`SELECT timescaledb_pre_restore();`)

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

	db.Exec(`SELECT timescaledb_post_restore();`)

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

// Delete ...
func (b *Backup) Delete(name string) (err error) {
	log.Infof("remove backup file %s", name)

	file := path.Join(b.cfg.Path, name)

	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		err = errors.Wrap(apperr.ErrBackupNotFound, fmt.Sprintf("path %s", file))
		return
	}

	if err = os.RemoveAll(file); err != nil {
		return
	}

	b.eventBus.Publish("system/services/backup", events.EventRemovedBackup{
		Name: name,
	})

	return
}

// UploadBackup ...
func (b *Backup) UploadBackup(ctx context.Context, reader *bufio.Reader, fileName string) (newFile *m.Backup, err error) {

	var list []*m.Backup
	if list, _, err = b.List(ctx, 999, 0, "", ""); err != nil {
		return
	}

	for _, file := range list {
		if fileName == file.Name {
			err = apperr.ErrBackupNameNotUnique
			return
		}
	}

	buffer := bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 128)

	var count int
	for {
		if count, err = reader.Read(part); err != nil {
			break
		}
		buffer.Write(part[:count])
	}
	if err != io.EOF {
	} else {
		err = nil
	}

	contentType := http.DetectContentType(buffer.Bytes())
	log.Infof("Content-type from buffer, %s", contentType)

	//create destination file making sure the path is writeable.
	var dst *os.File
	filePath := filepath.Join(b.cfg.Path, fileName)
	if dst, err = os.Create(filePath); err != nil {
		return
	}

	defer dst.Close()

	//copy the uploaded file to the destination file
	if _, err = io.Copy(dst, buffer); err != nil {
		return
	}

	size, _ := common.GetFileSize(filePath)
	newFile = &m.Backup{
		Name:     fileName,
		Size:     size,
		MimeType: contentType,
	}

	b.eventBus.Publish("system/services/backup", events.EventUploadedBackup{
		Name: fileName,
	})

	return
}
