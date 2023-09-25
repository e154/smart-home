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

package logging_ws

import (
	"context"
	"encoding/json"

	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/stream"
	"go.uber.org/fx"
)

var (
	log = logger.MustGetLogger("logging_ws")
)

type LoggingWs struct {
	logging *logging.Logging
	stream  *stream.Stream
}

func NewLogWsSaver(lc fx.Lifecycle,
	logging *logging.Logging,
	stream *stream.Stream) *LoggingWs {

	l := &LoggingWs{
		logging: logging,
		stream:  stream,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return l.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return l.Shutdown(ctx)
		},
	})

	return l
}

func (l *LoggingWs) Start(ctx context.Context) (err error) {
	log.Info("start")
	l.logging.SetWsSaver(l)
	return
}

func (l *LoggingWs) Shutdown(ctx context.Context) (err error) {
	log.Info("shutdown")
	l.logging.SetWsSaver(nil)
	return
}

// Save ...
func (l *LoggingWs) Save(log m.Log) {
	b, _ := json.Marshal(log)
	l.stream.Broadcast("log", b)
}
