package server

import "github.com/op/go-logging"

type ServerLogger struct {
	Logger *logging.Logger
}

func (s ServerLogger) Write(b []byte) (i int, err error) {
	s.Logger.Info(string(b))
	return
}
