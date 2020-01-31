package gate_client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type WsClient struct {
	settings   Settings
	delta      time.Duration
	interrupt  chan struct{}
	quitWorker chan struct{}
	cb         IWsCallback
	reConnect  bool
	mx         *sync.Mutex
	status     string
	counter    int
	sync.Mutex
	conn *websocket.Conn
}

func NewWsClient(
	cb IWsCallback) *WsClient {
	client := &WsClient{
		interrupt:  make(chan struct{}),
		quitWorker: make(chan struct{}),
		cb:         cb,
		mx:         &sync.Mutex{},
	}

	go func() {
		for {
			client.mx.Lock()
			if client.status == GateStatusQuit {
				client.mx.Unlock()
				return
			}
			client.mx.Unlock()

			client.connect()
			time.Sleep(time.Second)
		}
	}()
	return client
}

func (client *WsClient) UpdateSettings(settings Settings) {
	client.mx.Lock()
	defer client.mx.Unlock()

	client.settings = settings
	client.reConnect = true
}

func (client *WsClient) connect() {

	client.mx.Lock()
	client.status = GateStatusWait

	if !client.settings.Valid() {
		client.mx.Unlock()
		return
	}

	if !client.settings.Enabled {
		client.mx.Unlock()
		return
	}
	client.mx.Unlock()

	var err error
	ticker := time.NewTicker(pingPeriod)
	client.counter++
	defer func() {
		client.counter--

		client.mx.Lock()
		client.status = GateStatusNotConnected
		client.reConnect = false
		client.mx.Unlock()

		ticker.Stop()

		go client.cb.onClosed()

		if client.conn != nil {
			_ = client.conn.Close()
		}

		//if err != nil {
		//	if strings.Contains(err.Error(), "connection refused") {
		//		return
		//	}
		//	if strings.Contains(err.Error(), "bad handshake") {
		//		return
		//	}
		//	if strings.Contains(err.Error(), "use of closed network connection") {
		//		return
		//	}
		//	log.Debug(err.Error())
		//}
	}()

	var uri *url.URL
	if uri, err = url.Parse(client.settings.Address); err != nil {
		return
	}

	uri.Path = "ws"

	if uri.Scheme == "http" {
		uri.Scheme = "ws"
	} else {
		uri.Scheme = "wss"
	}

	requestHeader := http.Header{
		"X-Client-Type": {ClientTypeServer},
	}
	if client.settings.GateServerToken != "" {
		requestHeader.Add("X-API-Key", client.settings.GateServerToken)
		//log.Debugf("X-API-Key: %v", client.settings.GateServerToken)
	}

	if client.conn, _, err = websocket.DefaultDialer.Dial(uri.String(), requestHeader); err != nil {
		return
	}

	log.Info("gate connected ...")
	client.mx.Lock()
	client.status = GateStatusConnected
	client.mx.Unlock()

	loseChan := make(chan struct{})
	var messageType int
	var message []byte

	client.conn.SetCloseHandler(func(code int, text string) error {
		log.Warning("connection closed")

		loseChan <- struct{}{}
		return nil
	})

	go client.cb.onConnected()

	if err = client.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}

	client.conn.SetPongHandler(func(string) error {
		_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go func() {
		defer close(loseChan)
		for {

			if messageType, message, err = client.conn.ReadMessage(); err != nil {
				//log.Error(err.Error())
				loseChan <- struct{}{}
				break
			}
			switch messageType {
			case websocket.TextMessage:
				//fmt.Printf("recv: %s\n", string(message))
				go client.cb.onMessage(message)
			default:
				log.Warningf("unknown message type(%v)", messageType)
			}
		}
	}()

	for {
		select {
		case <-ticker.C:
			client.mx.Lock()
			if client.reConnect {
				//_ = client.write(websocket.CloseMessage, []byte{})
				fmt.Println("reconnect...")
				client.mx.Unlock()
				return
			}
			client.mx.Unlock()

			if err := client.write(websocket.PingMessage, []byte{}); err != nil {
				log.Error(err.Error())
				return
			}
		case <-client.interrupt:
		case <-loseChan:
			return
		}
	}
}

func (client *WsClient) Close() {

	client.mx.Lock()
	defer client.mx.Unlock()

	if client.status == GateStatusQuit {
		return
	}
	client.status = GateStatusQuit
	if client.status == GateStatusConnected {
		client.interrupt <- struct{}{}
	}
}

func (client *WsClient) write(opCode int, payload []byte) (err error) {

	client.mx.Lock()
	if client.status != GateStatusConnected {
		client.mx.Unlock()
		return
	}
	client.mx.Unlock()

	client.Lock()
	defer client.Unlock()

	if err = client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return
	}

	err = client.conn.WriteMessage(opCode, payload)

	return
}

func (client *WsClient) Status() string {
	client.mx.Lock()
	defer client.mx.Unlock()

	return client.status
}
