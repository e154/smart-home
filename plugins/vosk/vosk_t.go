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

//go:build test
// +build test

package vosk

import (
	"github.com/e154/smart-home/common/web"
	"io"
)

type Vosk struct{}

func (v Vosk) Start() {
}

func (v Vosk) Shutdown() {
}

func (v Vosk) LoadModel() error {
	return nil
}

func (v Vosk) STT(reader io.Reader, withGrm bool) (string, error) {
	return "", nil
}

func (v Vosk) DownloadModel(language string) error {
	return nil
}

func NewVosk(modelPath, sttLanguage string, crawler web.Crawler) *Vosk {
	return &Vosk{}
}
