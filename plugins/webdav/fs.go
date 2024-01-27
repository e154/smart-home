// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

package webdav

import (
	"context"
	"github.com/pkg/errors"
	"os"

	"github.com/spf13/afero"
	"golang.org/x/net/webdav"
)

type FS struct {
	afero.Fs
	onRemoveHandler func(context.Context, string) error
}

func NewFS(onRemoveHandler func(context.Context, string) error) *FS {
	return &FS{
		onRemoveHandler: onRemoveHandler,
		Fs:              afero.NewMemMapFs(),
	}
}

func (f *FS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return errors.New("operation not allowed")
}

func (f *FS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	_, err := f.Fs.Stat(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Infof("created new file %s", name)
		} else {
			return nil, err
		}
	}
	return f.Fs.OpenFile(name, flag, perm)
}

func (f *FS) RemoveAll(ctx context.Context, name string) error {
	if err := f.onRemoveHandler(ctx, name); err != nil {
		return err
	}
	return f.Fs.RemoveAll(name)
}

func (f *FS) Rename(ctx context.Context, oldName, newName string) error {
	return f.Fs.Rename(oldName, newName)
}

func (f *FS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	fileInfo, err := f.Fs.Stat(name)
	if err != nil {
		return nil, err
	}
	return fileInfo, err
}
