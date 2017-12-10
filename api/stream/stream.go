package stream

import (
	"net/http"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego"
)

type ConnectType string

const (
	SOCKJS ConnectType = ConnectType("sockjs")
	WEBSOCK ConnectType = ConnectType("websock")
)

var (
	h http.Handler
)

type StreamCotroller struct {
	beego.Controller
}

func (w *StreamCotroller) SockJs() {
	h = sockjs.NewHandler("/sockjs", sockjs.DefaultOptions, w.echoHandler)

	// CORS
	w.Ctx.ResponseWriter.Header().Del("Access-Control-Allow-Credentials")
	h.ServeHTTP(w.Ctx.ResponseWriter, w.Ctx.Request)
}

func (w *StreamCotroller) Ws() {

	// CORS
	w.Ctx.ResponseWriter.Header().Del("Access-Control-Allow-Credentials")

	conn, err := websocket.Upgrade(w.Ctx.ResponseWriter, w.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	hub := GetHub()
	client := &Client{
		ConnType: WEBSOCK,
		Connect: conn,
		Ip:w.Ctx.Input.IP(),
		Referer: w.Ctx.Input.Header("Referer"),
		UserAgent: w.Ctx.Request.UserAgent(),
		Send: make(chan []byte),
	}

	go client.WritePump()
	hub.AddClient(client)
}

func (w *StreamCotroller) echoHandler(session sockjs.Session) {

	hub := GetHub()
	client := &Client{
		ConnType: SOCKJS,
		Session:session,
		Ip:w.Ctx.Input.IP(),
		Referer: w.Ctx.Input.Header("Referer"),
		UserAgent: w.Ctx.Request.UserAgent(),
		Send: make(chan []byte),
	}

	hub.AddClient(client)
}

func (c *StreamCotroller) ErrHan(code int, message string) {

	switch code {
	case 401, 403:

	}

	c.Ctx.ResponseWriter.WriteHeader(code)
	c.Data["json"] = &map[string]interface{}{"status":"error", "message": message}
	c.ServeJSON()
}
