package logging

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/graceful_service"
	"time"
)

type LogDbSaver struct {
	adaptors  *adaptors.Adaptors
	isStarted bool
	pool      chan m.Log
	quit      chan struct{}
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

	if l.isStarted {
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

	l.isStarted = true
}

func (l *LogDbSaver) Shutdown() {
	if !l.isStarted {
		return
	}
	l.isStarted = false
	l.quit <- struct{}{}
	close(l.quit)
	close(l.pool)
}

func (l *LogDbSaver) Save(log m.Log) {
	if !l.isStarted {
		return
	}
	l.pool <- log
}
