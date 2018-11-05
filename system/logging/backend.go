package logging

import (
	"github.com/op/go-logging"
	"github.com/sirupsen/logrus"
)

type LogBackend struct {
	L *logrus.Logger
}

func NewLogBackend(logger *logrus.Logger) *LogBackend {
	return &LogBackend{L: logger}
}

func (b *LogBackend) Log(level logging.Level, calldepth int, rec *logging.Record) error {

	s := rec.Formatted(calldepth + 1)
	switch level {
	case logging.CRITICAL:
		b.L.Level = logrus.FatalLevel
		b.L.Fatal(s)
	case logging.ERROR:
		b.L.Level = logrus.ErrorLevel
		b.L.Error(s)
	case logging.WARNING:
		b.L.Level = logrus.WarnLevel
		b.L.Warning(s)
	case logging.INFO, logging.NOTICE:
		b.L.Level = logrus.InfoLevel
		b.L.Info(s)
	case logging.DEBUG:
		b.L.Level = logrus.DebugLevel
		b.L.Debug(s)
	}
	return nil
}
