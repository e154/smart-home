package stream

import (
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("stream")
)

type StreamService struct {
	Hub *Hub
}

func NewStreamService(hub *Hub) *StreamService {
	return &StreamService{
		Hub: hub,
	}
}

func (s *StreamService) Broadcast(message []byte) {
	s.Hub.Broadcast(message)
}

func (s *StreamService) Subscribe(command string, f func(client *Client, value interface{})) {
	s.Hub.Subscribe(command, f)
}

func (s *StreamService) UnSubscribe(command string) {
	s.Hub.UnSubscribe(command)
}

func (s *StreamService) AddClient(client *Client) {
	s.Hub.AddClient(client)
}
