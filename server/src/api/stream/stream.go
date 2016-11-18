package stream

import (
	"net/http"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"github.com/astaxie/beego"
)

var (
	h http.Handler
)

type StreamCotroller struct {
	beego.Controller
}

func (w *StreamCotroller) Get() {
	h = sockjs.NewHandler("/ws", sockjs.DefaultOptions, w.echoHandler)
	h.ServeHTTP(w.Ctx.ResponseWriter, w.Ctx.Request)
}

func (w *StreamCotroller) echoHandler(session sockjs.Session) {

	hub := GetHub()
	client := &Client{
		Session:session,
		Ip:w.Ctx.Input.IP(),
		Referer: w.Ctx.Input.Header("Referer"),
		UserAgent: w.Ctx.Request.UserAgent(),
	}
	//w.SetSession("clientinfo", client)
	hub.AddClient(client)
}
