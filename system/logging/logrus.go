package logging

import "github.com/sirupsen/logrus"

func NewLogrus() *logrus.Logger {
	return logrus.New()
}