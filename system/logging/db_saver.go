// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

package logging

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/graceful_service"
	"sync"
	"time"
)

type LogDbSaver struct {
	adaptors *adaptors.Adaptors
	pool     chan m.Log
	quit     chan struct{}
	sync.Mutex
	isRunning bool
}

func NewLogDbSaver(adaptors *adaptors.Adaptors,
	graceful *graceful_service.GracefulService, ) *LogDbSaver {
	saver := &LogDbSaver{
		adaptors: adaptors,
		pool:     make(chan m.Log),
		quit:     make(chan struct{}),
	}

	graceful.Subscribe(saver)

	saver.Start()

	return saver
}

func (l *LogDbSaver) Start() {

	if l.safeIsRunning() {
		return
	}

	go func() {

		logList := make([]*m.Log, 0, 50)
		ticker := time.NewTicker(time.Second * 5)
		defer func() {
			ticker.Stop()
		}()

		update := func() {
			_ = l.adaptors.Log.AddMultiple(logList)
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

	l.safeSetIsRunning(true)
}

func (l *LogDbSaver) Shutdown() {
	if !l.safeIsRunning() {
		return
	}
	l.safeSetIsRunning(false)
	l.quit <- struct{}{}
	close(l.quit)
	close(l.pool)
}

func (l *LogDbSaver) Save(log m.Log) {
	if !l.safeIsRunning() {
		return
	}
	l.pool <- log
}

func (l *LogDbSaver) safeIsRunning() bool {
	l.Lock()
	defer l.Unlock()
	return l.isRunning
}

func (l *LogDbSaver) safeSetIsRunning(v bool) {
	l.Lock()
	l.isRunning = v
	l.Unlock()
}
