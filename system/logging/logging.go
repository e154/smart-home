package logging

import (
	"os"
	"github.com/sirupsen/logrus"
	"github.com/op/go-logging"
	"github.com/olivere/elastic"
)

var (
	client *elastic.Client
	log    = logrus.New()
	format = logging.MustStringFormatter(
		"SRT.%{module}.%{shortfile}.%{shortfunc}() > %{message}",
	)
)

func Initialize(log *logrus.Logger) {

	log.Out = os.Stdout

	log1 := NewLogBackend(log)
	log1F := logging.NewBackendFormatter(log1, format)
	logging.SetBackend(log1F)
}

func NewLogrus() *logrus.Logger {
	return logrus.New()
}