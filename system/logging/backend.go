package logging

import (
	"os"
	"github.com/op/go-logging"
	"github.com/sirupsen/logrus"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/common"
)

type LogBackend struct {
	L        *logrus.Logger
	adaptors *adaptors.Adaptors
	oldLog   *m.Log
}

func NewLogBackend(
	logger *logrus.Logger,
	adaptors *adaptors.Adaptors) (back *LogBackend) {
	back = &LogBackend{
		L:        logger,
		adaptors: adaptors,
	}
	back.L.Out = os.Stdout
	return
}

func (b *LogBackend) Log(level logging.Level, calldepth int, rec *logging.Record) (err error) {

	var logLevel common.LogLevel
	s := rec.Formatted(calldepth + 1)
	switch level {
	case logging.CRITICAL:
		b.L.Level = logrus.FatalLevel
		b.L.Fatal(s)
		logLevel = "Critical"
	case logging.ERROR:
		b.L.Level = logrus.ErrorLevel
		b.L.Error(s)
		logLevel = "Error"
	case logging.WARNING:
		b.L.Level = logrus.WarnLevel
		b.L.Warning(s)
		logLevel = "Warning"
	case logging.INFO, logging.NOTICE:
		b.L.Level = logrus.InfoLevel
		b.L.Info(s)
		logLevel = "Info"
	case logging.DEBUG:
		b.L.Level = logrus.DebugLevel
		b.L.Debug(s)
		logLevel = "Debug"
	}

	//TODO optimise
	if b.oldLog != nil {
		if b.oldLog.Body == rec.Message() && b.oldLog.Level == logLevel {
			return
		}
	}

	record := &m.Log{
		Level:     logLevel,
		Body:      rec.Message(),
		CreatedAt: rec.Time,
	}

	b.oldLog = record

	_, err = b.adaptors.Log.Add(record)

	return nil
}
