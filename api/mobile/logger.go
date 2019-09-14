package mobile

import "github.com/op/go-logging"

type MobileServerLogger struct {
	Logger *logging.Logger
}

func (s MobileServerLogger) Write(b []byte) (i int, err error) {
	s.Logger.Info(string(b))
	return
}
