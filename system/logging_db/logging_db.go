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

package logging_db

import (
	"context"
	"time"

	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/system/logging"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"go.uber.org/atomic"
	"go.uber.org/fx"
)

var (
	log = logger.MustGetLogger("logging_db")
)

// LogDbSaver ...
type LogDbSaver struct {
	adaptors  *adaptors.Adaptors
	pool      chan m.Log
	quit      chan struct{}
	isRunning *atomic.Bool
}

// NewLogDbSaver ...
func NewLogDbSaver(lc fx.Lifecycle,
	adaptors *adaptors.Adaptors) logging.ISaver {
	saver := &LogDbSaver{
		adaptors:  adaptors,
		pool:      make(chan m.Log),
		quit:      make(chan struct{}),
		isRunning: atomic.NewBool(false),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			saver.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			saver.Shutdown()
			return nil
		},
	})

	return saver
}

// Start ...
func (l *LogDbSaver) Start() {
	log.Info("start")

	if l.isRunning.Load() {
		return
	}

	go func() {

		logList := make([]*m.Log, 0, 50)
		ticker := time.NewTicker(time.Second * 5)
		defer func() {
			ticker.Stop()
		}()

		update := func() {
			_ = l.adaptors.Log.AddMultiple(context.Background(), logList)
			logList = make([]*m.Log, 0, 50)
		}

		for {
			select {
			case <-ticker.C:
				if len(logList) > 0 {
					update()
				}
			case logMsg := <-l.pool:
				logList = append(logList, &logMsg)
				if len(logList) >= 50 {
					update()
				}
			case <-l.quit:
				return
			}
		}
	}()

	l.isRunning.Store(true)
}

// Shutdown ...
func (l *LogDbSaver) Shutdown() {
	log.Info("shutdown")

	if !l.isRunning.Load() {
		return
	}
	l.isRunning.Store(false)
	l.quit <- struct{}{}
	close(l.quit)
	close(l.pool)
}

// Save ...
func (l *LogDbSaver) Save(log m.Log) {
	if !l.isRunning.Load() {
		return
	}
	l.pool <- log //todo fix: some time panic: send on closed channel
}
