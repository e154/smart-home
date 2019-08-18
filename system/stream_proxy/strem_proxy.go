package stream_proxy

import (
	"fmt"
	"github.com/e154/smart-home/api/server"
	"github.com/e154/smart-home/common/debug"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/stream"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"net/http"
	"net/http/httptest"
)

var (
	log = logging.MustGetLogger("stream_proxy")
)

type StreamProxy struct {
	engine        *gin.Engine
	streamService *stream.StreamService
}

func NewStreamProxy(httpsServer *server.Server,
	streamService *stream.StreamService,
	graceful *graceful_service.GracefulService) (proxy *StreamProxy) {
	proxy = &StreamProxy{
		engine:        httpsServer.GetEngine(),
		streamService: streamService,
	}

	graceful.Subscribe(proxy)

	return
}

func (s *StreamProxy) Start() {
	log.Info("start stream proxy")
	s.streamService.Subscribe("chanel.server", s.DoAction)
}

func (s *StreamProxy) Shutdown() {
	s.streamService.UnSubscribe("chanel.server")
	return
}

func (s *StreamProxy) DoAction(client *stream.Client, message stream.Message) {

	fmt.Println("------")
	debug.Println(client)
	fmt.Println("------")
	debug.Println(message)

	return
}

func (s *StreamProxy) execRequest() {

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v1/", nil)
	request.SetBasicAuth("admin", "admin")

	s.engine.ServeHTTP(recorder, request)
	fmt.Println(recorder.Code)
	fmt.Println(recorder.Body)
}

