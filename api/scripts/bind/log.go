package bind

import "github.com/e154/smart-home/api/log"

type LogBind struct {}

func (b *LogBind) Info(v ...interface{}) { log.Info(v) }
func (b *LogBind) Warn(v ...interface{}) { log.Warn(v) }
func (b *LogBind) Debug(v ...interface{}) { log.Debug(v) }
func (b *LogBind) Error(v ...interface{}) { log.Error(v) }