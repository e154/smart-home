package logging

import (
	"github.com/op/go-logging"
)

var (
	format = logging.MustStringFormatter(
		"SRT.%{module}.%{shortfile}.%{shortfunc}() > %{message}",
	)
)

func Initialize(log1 *LogBackend) {
	log1F := logging.NewBackendFormatter(log1, format)
	logging.SetBackend(log1F)
}
