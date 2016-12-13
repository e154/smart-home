package log

import (
	"time"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"runtime"
	"fmt"
	"../models"
)

// brush is a color join function
type brush func(string) string

// consoleWriter implements LoggerInterface and writes messages to terminal.
type logger struct {
	Level    int  `json:"level"`
	Colorful bool `json:"color"`
}

// newBrush return a fix color Brush
func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []brush{
	newBrush("1;37"), // Emergency          white
	newBrush("1;36"), // Alert              cyan
	newBrush("1;35"), // Critical           magenta
	newBrush("1;31"), // Error              red
	newBrush("1;33"), // Warning            yellow
	newBrush("1;32"), // Notice             green
	newBrush("1;34"), // Informational      blue
	newBrush("1;34"), // Debug              blue
}

// NewConsole create ConsoleWriter returning as LoggerInterface.
func SmartLogger() logs.Logger {
	cw := &logger{
		Colorful: runtime.GOOS != "windows",
	}

	return cw
}

// jsonConfig like '{"level":LevelTrace}'.
func (c *logger) Init(jsonConfig string) error {

	return nil
}

// WriteMsg write message in console.
func (c *logger) WriteMsg(when time.Time, msg string, level int) (err error) {

	//...
	if err = c.story(when, msg, level); err != nil {
		return err
	}

	if(beego.BConfig.RunMode != "dev") {
		return
	}

	if c.Colorful {
		msg = colors[level](msg)
	}

	fmt.Println(when.Format("2006/01/02 15:04:05"), msg)

	return
}

// Destroy implementing method. empty.
func (c *logger) Destroy() {

}

// Flush implementing method. empty.
func (c *logger) Flush() {

}

func (c *logger) story(when time.Time, msg string, l int) (err error) {

	var level string
	switch l {
	case logs.LevelEmergency:
		level = "Emergency"
	case logs.LevelAlert:
		level = "Alert"
	case logs.LevelCritical:
		level = "Critical"
	case logs.LevelError:
		level = "Error"
	case logs.LevelWarning:
		level = "Warning"
	case logs.LevelNotice:
		level = "Notice"
	case logs.LevelInformational:
		level = "Info"
	case logs.LevelDebug:
		level = "Debug"
	}

	log := &models.Log{
		Level: level,
		Created_at: when,
		Body: msg,
	}

	if _, err = models.AddLog(log); err != nil {
		beego.Error("error", err.Error())
	}

	return
}