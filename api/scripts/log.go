package scripts

import (
	"github.com/e154/smart-home/api/log"
)

type Log struct {

}

func (l *Log) Info(v ...interface{}) {
	log.Info(v)
}

func (l *Log) Warn(v ...interface{}) {
	log.Warn(v)
}

func (l *Log) Debug(v ...interface{}) {
	log.Debug(v)
}

func (l *Log) Error(v ...interface{}) {
	log.Error(v)
}
