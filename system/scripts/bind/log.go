package bind

import "github.com/op/go-logging"

var (
	log = logging.MustGetLogger("js")
)

// Javascript Binding
//
// IC.Log
// 	 .info()
// 	 .warn()
// 	 .error()
// 	 .debug()
//
type LogBind struct{}

func (b *LogBind) Info(v ...interface{})  { log.Infof("%v", v) }
func (b *LogBind) Warn(v ...interface{})  { log.Warningf("%v", v) }
func (b *LogBind) Debug(v ...interface{}) { log.Debugf("%v", v) }
func (b *LogBind) Error(v ...interface{}) { log.Errorf("%v", v) }
