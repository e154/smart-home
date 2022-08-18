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

package backup

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"go.uber.org/fx"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	app "github.com/e154/smart-home/common/app"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
)

var (
	log = logger.MustGetLogger("backup")
)

// Backup ...
type Backup struct {
	cfg          *BackupConfig
	Options      []string
	db           *gorm.DB
	restoreImage string
}

// NewBackup ...
func NewBackup(lc fx.Lifecycle,
	cfg *BackupConfig,
	db *gorm.DB) *Backup {

	backup := &Backup{
		cfg: cfg,
		db:  db,
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

	options := b.dumpOptions()

	tmpDir := path.Join(os.TempDir(), "smart_home")
	if err = os.MkdirAll(tmpDir, 0755); err != nil {
		return
	}

	// filename
	filename := path.Join(tmpDir, "database.sql")
	options = append(options, "-f", filename)

	cmd := exec.Command("pg_dump", options...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", b.cfg.PgPass))

	log.Infof("run command %s", cmd.String())

	_, err = cmd.CombinedOutput()
	if err != nil {
		return
	}

	err = zipit([]string{path.Join("data", "file_storage"), filename}, path.Join(b.cfg.Path, fmt.Sprintf("%s.zip", time.Now().Format("2006-01-02T15:04:05.999"))))
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

	log.Info("Purge database")

	if err = b.db.Exec(`DROP SCHEMA IF EXISTS "public" CASCADE;`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed exec sql command"), err.Error())
		return
	}
	if err = b.db.Exec(`CREATE SCHEMA "public";`).Error; err != nil {
		err = errors.Wrap(fmt.Errorf("failed exec sql command"), err.Error())
		return
	}

	options := b.restoreOptions()

	options = append(options, "-f", path.Join(tmpDir, "database.sql"))

	cmd := exec.Command("psql", options...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", b.cfg.PgPass))

	log.Infof("command: %s", cmd.String())

	if _, err = cmd.CombinedOutput(); err != nil {
		err = errors.Wrap(fmt.Errorf("failed combine command"), err.Error())
		return
	}

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

	log.Info("complete")

	return
}

func (b Backup) dumpOptions() []string {
	options := b.Options

	// db name
	if b.cfg.PgName != "" {
		options = append(options, "-d")
		options = append(options, b.cfg.PgName)
	}

	// host
	if b.cfg.PgHost != "" {
		options = append(options, "-h")
		options = append(options, b.cfg.PgHost)
	}

	// port
	if b.cfg.PgPort != "" {
		options = append(options, "-p")
		options = append(options, b.cfg.PgPort)
	}

	// user
	if b.cfg.PgUser != "" {
		options = append(options, "-U")
		options = append(options, b.cfg.PgUser)
	}

	// compress level
	//options = append(options, "-Z", "9")

	// formats
	options = append(options, "-F", "p")

	return options
}

func (b Backup) restoreOptions() []string {
	options := b.Options

	// db name
	if b.cfg.PgName != "" {
		options = append(options, "-d")
		options = append(options, b.cfg.PgName)
	}

	// host
	if b.cfg.PgHost != "" {
		options = append(options, "-h")
		options = append(options, b.cfg.PgHost)
	}

	// port
	if b.cfg.PgPort != "" {
		options = append(options, "-p")
		options = append(options, b.cfg.PgPort)
	}

	// user
	if b.cfg.PgUser != "" {
		options = append(options, "-U")
		options = append(options, b.cfg.PgUser)
	}

	return options
}
