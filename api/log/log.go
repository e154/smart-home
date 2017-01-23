package log

import (
	"strings"
	"github.com/astaxie/beego/logs"
)

// BeeLogger references the used application logger.
var log *logs.BeeLogger

// Emergency log a message at emergency level.
func Emergency(v ...interface{}) {
	log.Emergency(generateFmtStr(len(v)), v...)
}

func Emergencyf(a string, v ...interface{}) {
	log.Emergency(a, v...)
}


// Alert log a message at alert level.
func Alert(v ...interface{}) {
	log.Alert(generateFmtStr(len(v)), v...)
}

func Alertf(a string, v ...interface{}) {
	log.Alert(a, v...)
}


// Critical log a message at critical level.
func Critical(v ...interface{}) {
	log.Critical(generateFmtStr(len(v)), v...)
}

func Criticalf(a string, v ...interface{}) {
	log.Critical(a, v...)
}

// Error log a message at error level.
func Error(v ...interface{}) {
	log.Error(generateFmtStr(len(v)), v...)
}

func Errorf(a string, v ...interface{}) {
	log.Error(a, v...)
}

// Warn compatibility alias for Warning()
func Warn(v ...interface{}) {
	log.Warn(generateFmtStr(len(v)), v...)
}

func Warnf(a string, v ...interface{}) {
	log.Warn(a, v...)
}

// Notice log a message at notice level.
func Notice(v ...interface{}) {
	log.Notice(generateFmtStr(len(v)), v...)
}

func Noticef(a string, v ...interface{}) {
	log.Notice(a, v...)
}

// Info compatibility alias for Warning()
func Info(v ...interface{}) {
	log.Info(generateFmtStr(len(v)), v...)
}

func Infof(a string, v ...interface{}) {
	log.Info(a, v...)
}

// Debug log a message at debug level.
func Debug(v ...interface{}) {
	log.Debug(generateFmtStr(len(v)), v...)
}

func Debugf(a string, v ...interface{}) {
	log.Debug(a, v...)
}

// Trace log a message at trace level.
// compatibility alias for Warning()
func Trace(v ...interface{}) {
	log.Trace(generateFmtStr(len(v)), v...)
}

func Tracef(a string, v ...interface{}) {
	log.Trace(a, v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}

func init() {
	if log != nil {
		return
	}

	log = logs.NewLogger(10000)
	logs.Register("smart", SmartLogger)
	log.SetLogger("smart", "")
}