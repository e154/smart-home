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

package vosk

import (
	"io"
	"path"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	Name         = "vosk"
	Version      = "0.0.1"
	AttrModel    = "model"
	FunctionName = "automationTriggerStt"
	sampleRate   = 16000.0
	URLPrefix    = "https://alphacephei.com/vosk/models/"
	defaultModel = "vosk-model-small-en-us-0.15"
)

var (
	modelPath = path.Join("data", "vosk", "models")
)

// NewSettings ...
func NewSettings() map[string]*m.Attribute {
	return map[string]*m.Attribute{
		AttrModel: {
			Name:  AttrModel,
			Type:  common.AttributeString,
			Value: defaultModel,
		},
	}
}

func NewTriggerParams() m.TriggerParams {
	return m.TriggerParams{
		Script:     true,
		Required:   []string{},
		Attributes: m.Attributes{},
	}
}

type STT interface {
	Start()
	Shutdown()
	LoadModel() error
	STT(reader io.Reader, withGrm bool) (string, error)
	DownloadModel(language string) error
}
