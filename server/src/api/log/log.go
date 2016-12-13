package log

import (
	"strings"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
)

// BeeLogger references the used application logger.
var log *logs.BeeLogger

// SetLevel sets the global log level used by the simple logger.
func SetLevel(l int) {
	log.SetLevel(l)
}

// SetLogFuncCall set the CallDepth, default is 3
func SetLogFuncCall(b bool) {
	logs.SetLogFuncCall(b)
}

// SetLogger sets a new logger.
func SetLogger(adaptername string, config string) error {
	return log.SetLogger(adaptername, config)
}

// Emergency log a message at emergency level.
func Emergency(v ...interface{}) {
	log.Emergency(generateFmtStr(len(v)), v...)
}

// Alert log a message at alert level.
func Alert(v ...interface{}) {
	log.Alert(generateFmtStr(len(v)), v...)
}

// Critical log a message at critical level.
func Critical(v ...interface{}) {
	log.Critical(generateFmtStr(len(v)), v...)
}

// Error log a message at error level.
func Error(v ...interface{}) {
	log.Error(generateFmtStr(len(v)), v...)
}

// Warning log a message at warning level.
func Warning(v ...interface{}) {
	log.Warning(generateFmtStr(len(v)), v...)
}

// Warn compatibility alias for Warning()
func Warn(v ...interface{}) {
	log.Warn(generateFmtStr(len(v)), v...)
}

// Notice log a message at notice level.
func Notice(v ...interface{}) {
	log.Notice(generateFmtStr(len(v)), v...)
}

// Informational log a message at info level.
func Informational(v ...interface{}) {
	log.Informational(generateFmtStr(len(v)), v...)
}

// Info compatibility alias for Warning()
func Info(v ...interface{}) {
	beego.Info()
	log.Info(generateFmtStr(len(v)), v...)
}

// Debug log a message at debug level.
func Debug(v ...interface{}) {
	log.Debug(generateFmtStr(len(v)), v...)
}

// Trace log a message at trace level.
// compatibility alias for Warning()
func Trace(v ...interface{}) {
	log.Trace(generateFmtStr(len(v)), v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}

func Println(v ...interface{}) {
	Info(v...)
}

func Fatal(v ...interface{}) {
	Critical(v...)
}

func init() {
	if log != nil {
		return
	}

	log = logs.NewLogger(10000)
	logs.Register("smart", SmartLogger)
	log.SetLogger("smart", "")
}