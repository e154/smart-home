// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path"
)

type Local struct {
	cfg *Config
}

func NewLocal(cfg *Config) *Local {
	return &Local{cfg: cfg}
}

func (l *Local) New(tmpDir string) (err error) {

	options := l.dumpOptions()

	// get data
	var dataFile = path.Join(tmpDir, "data.sql")
	cmd := exec.Command("pg_dump", append(options, "-a", "-f", dataFile)...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", l.cfg.Password))

	log.Infof("run command: %s", cmd.String())

	if _, err = cmd.CombinedOutput(); err != nil {
		return err
	}

	// get scheme
	var schemeFile = path.Join(tmpDir, "scheme.sql")
	cmd = exec.Command("pg_dump", append(options, "-s", "-f", schemeFile)...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", l.cfg.Password))

	log.Infof("run command: %s", cmd.String())

	_, err = cmd.CombinedOutput()

	return
}

func (l *Local) Restore(path string) (err error) {

	options := l.restoreOptions()

	options = append(options, "-v", path)

	cmd := exec.Command("pg_restore", options...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", l.cfg.Password))

	log.Infof("command: %s", cmd.String())

	if _, err = cmd.CombinedOutput(); err != nil {
		err = errors.Wrap(fmt.Errorf("failed combine command"), err.Error())
	}

	return
}

func (l *Local) dumpOptions() []string {
	options := []string{}

	// db name
	if l.cfg.Name != "" {
		options = append(options, "-d", l.cfg.Name)
	}

	// host
	if l.cfg.Host != "" {
		options = append(options, "-h", l.cfg.Host)
	}

	// port
	if l.cfg.Port != "" {
		options = append(options, "-p", l.cfg.Port)
	}

	// user
	if l.cfg.User != "" {
		options = append(options, "-U", l.cfg.User)
	}

	// compress level
	//options = append(options, "-Z", "9")

	// formats
	options = append(options, "-F", "t")

	// etc
	options = append(options, "-w", "-v", "--disable-triggers", "--quote-all-identifiers", "-b", "-O", "-x")

	return options
}

func (l *Local) restoreOptions() []string {

	options := []string{}

	// db name
	if l.cfg.Name != "" {
		options = append(options, "-d", l.cfg.Name)
	}

	// host
	if l.cfg.Host != "" {
		options = append(options, "-h", l.cfg.Host)
	}

	// port
	if l.cfg.Port != "" {
		options = append(options, "-p", l.cfg.Port)
	}

	// user
	if l.cfg.User != "" {
		options = append(options, "-U", l.cfg.User)
	}

	return options
}
