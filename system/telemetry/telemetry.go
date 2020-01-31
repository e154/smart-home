package telemetry

import (
	"github.com/op/go-logging"
	"sync"
)

var (
	log = logging.MustGetLogger("telemetry")
)

type Telemetry struct {
	sync.Mutex
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
	s.Lock()
	defer s.Unlock()
	if s.subscribers[command] != nil {
		delete(s.subscribers, command)
	}
	s.subscribers[command] = f
}

func (s *Telemetry) UnSubscribe(command string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.subscribers[command]; ok {
		delete(s.subscribers, command)
	}
}

func (t *Telemetry) Broadcast(param interface{}) {
	t.Lock()
	defer t.Unlock()
	for _, f := range t.subscribers {
		f.Broadcast(param)
	}
}

func (t *Telemetry) BroadcastOne(param interface{}) {
	t.Lock()
	defer t.Unlock()
	for _, f := range t.subscribers {
		f.BroadcastOne(param)
	}
}
