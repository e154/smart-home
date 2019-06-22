package telemetry

import (
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("telemetry")
)

type Telemetry struct {
	subscribers map[string]ITelemetry
}

func NewTelemetry() (t2 *Telemetry, t1 ITelemetry) {
	t2 = &Telemetry{
		subscribers: make(map[string]ITelemetry),
	}
	t1 = t2
	return
}

func (s *Telemetry) Subscribe(command string, f ITelemetry) {
	log.Infof("subscribe %s", command)
	if s.subscribers[command] != nil {
		delete(s.subscribers, command)
	}
	s.subscribers[command] = f
}

func (s *Telemetry) UnSubscribe(command string) {
	if _, ok := s.subscribers[command]; ok {
		delete(s.subscribers, command)
	}
}

func (t *Telemetry) Broadcast(pack string) {
	for _, f := range t.subscribers {
		f.Broadcast(pack)
	}
}

func (t *Telemetry) BroadcastOne(pack string, deviceId int64, elementName string) {
	for _, f := range t.subscribers {
		f.BroadcastOne(pack, deviceId, elementName)
	}
}
