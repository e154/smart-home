package dashboard_models

import "github.com/op/go-logging"

var (
	log = logging.MustGetLogger("models")
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
