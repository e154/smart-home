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

package supervisor

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"plugin"
	"runtime"

	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/apperr"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

type ExternalPlugins struct {
	adaptors *adaptors.Adaptors
}

func NewExternalPlugins(adaptors *adaptors.Adaptors) *ExternalPlugins {
	return &ExternalPlugins{
		adaptors: adaptors,
	}
}

func (p *ExternalPlugins) uploadPlugin(ctx context.Context, reader *bufio.Reader) (plugin *m.Plugin, err error) {

	buffer := bytes.NewBuffer([]byte{})
	part := make([]byte, 128)

	var count int
	for {
		if count, err = reader.Read(part); err != nil {
			break
		}
		buffer.Write(part[:count])
	}

	if err != io.EOF {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrPluginUpload)
		return
	}

	contentType := http.DetectContentType(buffer.Bytes())

	switch contentType {
	case "application/x-gzip":
	default:
		log.Warnf("Unknown plugin archive type: %s", contentType)
	}

	var manifest *plugins.PluginManifest
	if manifest, err = p.readArchiveManifest(common.CopyBuffer(buffer)); err != nil {
		return
	}

	if manifest == nil {
		return nil, fmt.Errorf("manifest file not found or corrupted")
	}

	if manifest.OS != runtime.GOOS {
		return nil, fmt.Errorf("this plugin only for %s operating system, current operating system: %s", manifest.OS, runtime.GOOS)
	}

	if manifest.Arch != runtime.GOARCH {
		return nil, fmt.Errorf("this plugin only for %s architecture, current architecture: %s", manifest.Arch, runtime.GOARCH)
	}

	if err = p.checkArchive(common.CopyBuffer(buffer), manifest); err != nil {
		return nil, fmt.Errorf("%s: archive failed manifest check", err.Error())
	}

	plugin, err = p.adaptors.Plugin.GetByName(context.Background(), manifest.Name)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			err = nil
		} else {
			return nil, err
		}
	}

	pluginDir := filepath.Join("data", "plugins", manifest.Name)

	// mkdir
	if err = os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("%s: create dir %s failed", err.Error(), pluginDir)
	}

	// copy plugin.so
	to := filepath.Join(pluginDir, manifest.Plugin)
	if err = p.extractFileFromArchive(common.CopyBuffer(buffer), manifest.Plugin, to); err != nil {
		return nil, fmt.Errorf("%s: copy plugin file %s to %s failed", err.Error(), manifest.Plugin, to)
	}

	// copy libs
	for _, item := range manifest.Libs {
		to = filepath.Join(".", item)
		if err = p.extractFileFromArchive(common.CopyBuffer(buffer), item, to); err != nil {
			return nil, fmt.Errorf("%s: copy library file %s to %s failed", err.Error(), item, to)
		}
	}

	// copy asset
	for _, item := range manifest.Assets {
		to = filepath.Join(pluginDir, item)
		if err = p.extractFileFromArchive(common.CopyBuffer(buffer), item, to); err != nil {
			return nil, fmt.Errorf("%s: copy asset file %s to %s failed", err.Error(), item, to)
		}
	}

	// copy manifest.json
	to = filepath.Join(pluginDir, "manifest.json")
	if err = p.extractFileFromArchive(common.CopyBuffer(buffer), "manifest.json", to); err != nil {
		return nil, fmt.Errorf("%s: copy manifest file %s to %s failed", err.Error(), "manifest.json", to)
	}

	if plugin == nil {
		plugin = &m.Plugin{
			Name:     manifest.Name,
			Version:  manifest.Version,
			Settings: manifest.Settings,
			Enabled:  false,
			System:   false,
			Actor:    manifest.Actor,
			IsLoaded: false,
			Triggers: manifest.Triggers,
			External: true,
		}
	} else {
		plugin.Version = manifest.Version
		//plugin.Settings = manifest.Settings
		plugin.Triggers = manifest.Triggers
		plugin.Actor = manifest.Actor
		plugin.External = true
	}

	err = p.adaptors.Plugin.CreateOrUpdate(ctx, plugin)

	return
}

func (p *ExternalPlugins) removeExternalPlugin(ctx context.Context, pluginName string) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrPluginDelete)
		}
	}()

	var plugin *m.Plugin
	if plugin, err = p.adaptors.Plugin.GetByName(ctx, pluginName); err != nil {
		return
	}

	if !plugin.External {
		err = fmt.Errorf("%s: %w", "you can only remove external plugins", apperr.ErrPluginDelete)
		return
	}

	pluginDir := filepath.Join("data", "plugins", plugin.Name)

	// read manifest
	file, err := os.OpenFile(filepath.Join(pluginDir, "manifest.json"), os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("%s: can't open fole %s", err.Error(), filepath.Join(pluginDir, "manifest.json"))
	}
	defer file.Close()

	manifest, err := p.readManifest(file)
	if err != nil {
		return err
	}

	// check name
	if manifest.Name == "" || manifest.Name != plugin.Name {
		return fmt.Errorf("the name field is empty in the manifest file is incorrect")
	}

	if err = os.RemoveAll(pluginDir); err != nil {
		log.Warn(err.Error())
	}

	for _, item := range manifest.Libs {
		if err = os.RemoveAll(filepath.Join(".", item)); err != nil {
			log.Warn(err.Error())
		}
	}

	return p.adaptors.Plugin.Delete(ctx, plugin.Name)
}

func (p *ExternalPlugins) readManifest(r io.Reader) (*plugins.PluginManifest, error) {

	manifest := &plugins.PluginManifest{}
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%s: read manifest failed", err.Error())
	}

	if err = json.Unmarshal(data, manifest); err != nil {
		return nil, fmt.Errorf("%s: unmarshal manifest failed", err.Error())
	}

	return manifest, nil
}

func (p *ExternalPlugins) readArchiveManifest(buffer io.Reader) (*plugins.PluginManifest, error) {

	uncompressedStream, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("%s: extract tar failed", err.Error())
		}

		if header == nil {
			return nil, fmt.Errorf("%s: extract tar failed", "header is nil")
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			if header.FileInfo().Name() != "manifest.json" {
				continue
			}

			return p.readManifest(tarReader)
		}
	}

	return nil, nil
}

func (p *ExternalPlugins) extractFileFromArchive(buffer io.Reader, fileName, destName string) error {

	uncompressedStream, err := gzip.NewReader(buffer)
	if err != nil {
		return err
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("%s: extract tar failed", err.Error())
		}

		if header == nil {
			return fmt.Errorf("%s: extract tar failed", "header is nil")
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			if header.FileInfo().Name() != fileName {
				continue
			}

			dest, err := os.OpenFile(destName, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				return fmt.Errorf("%s: can't create file %s", err.Error(), destName)
			}
			defer dest.Close()

			_, err = io.Copy(dest, tarReader)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *ExternalPlugins) checkArchive(buffer io.Reader, manifest *plugins.PluginManifest) error {

	uncompressedStream, err := gzip.NewReader(buffer)
	if err != nil {
		return err
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)

	files := map[string]struct{}{}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("%s: extract tar failed", err.Error())
		}

		if header == nil {
			return fmt.Errorf("%s: extract tar failed", "header is nil")
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			files[header.FileInfo().Name()] = struct{}{}
		}
	}

	if manifest.Name == "" {
		return fmt.Errorf("name is empty")
	}

	if manifest.Description == "" {
		return fmt.Errorf("description is empty")
	}

	if manifest.Version == "" {
		return fmt.Errorf("version is empty")
	}

	if manifest.Plugin == "" {
		return fmt.Errorf("plugin is empty")
	}

	if _, ok := files[manifest.Plugin]; !ok {
		return fmt.Errorf("plugin not found")
	}

	for _, item := range manifest.Libs {
		if _, ok := files[manifest.Plugin]; !ok {
			return fmt.Errorf("lib %s not found", item)
		}
	}

	for _, item := range manifest.Assets {
		if _, ok := files[manifest.Plugin]; !ok {
			return fmt.Errorf("asset %s not found", item)
		}
	}

	return nil
}

func (p *ExternalPlugins) loadGoPlugin(pluginName string) error {

	if _, ok := pluginList.Load(pluginName); ok {
		return nil
	}

	dir := path.Join(pluginsDir, pluginName, "plugin.so")
	log.Infof("load external plugin %s", dir)
	plugin, err := plugin.Open(dir)
	if err != nil {
		return err
	}

	newFunc, err := plugin.Lookup("New")
	if err != nil {
		return err
	}

	pluggable, ok := newFunc.(func() plugins.Pluggable)
	if !ok {
		return errors.New("unexpected type from module symbol")
	}

	RegisterPlugin(pluginName, pluggable)

	return nil
}

func (p *ExternalPlugins) loadExternalPlugins() {

	var list plugins.PluginFileInfos

	_ = filepath.Walk(pluginsDir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() && info.Name() == "plugins" {
			return nil
		}
		if info.IsDir() {
			list = append(list, &plugins.PluginFileInfo{
				Name:     info.Name(),
				Size:     info.Size(),
				FileMode: info.Mode(),
				ModTime:  info.ModTime(),
				IsDir:    info.IsDir(),
			})

			return filepath.SkipDir
		}
		return nil
	})

	for _, item := range list {
		if !item.IsDir {
			continue
		}
		if err := p.loadGoPlugin(item.Name); err != nil {
			err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrPluginLoadExternal)
			log.Warn(err.Error())
		}
	}
}
