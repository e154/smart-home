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

package events

import (
	"github.com/e154/smart-home/pkg/common"
)

type EventCreatedBackup struct {
	Name      string `json:"name"`
	Scheduler bool   `json:"scheduler"`
}

type EventRemovedBackup struct {
	Name string `json:"name"`
}

type EventUploadedBackup struct {
	Name string `json:"name"`
}

type EventStartedRestore struct {
	Name string `json:"name"`
}

type CommandCreateBackup struct {
	Scheduler bool `json:"scheduler"`
}

type CommandClearStorage struct {
	Num int64 `json:"num"`
}

type CommandSendFileToTelegram struct {
	Filename  string          `json:"filename"`
	EntityId  common.EntityId `json:"entity_id"`
	Chunks    bool            `json:"chunks"`
	ChunkSize int             `json:"chunk_size"`
}
