// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package media

import (
	"net/http"
	"time"

	"github.com/deepch/vdk/format/mp4f"
	"github.com/e154/smart-home/internal/plugins/media/server"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// ControllerMedia ...
type ControllerMedia struct {
}

// NewControllerMedia ...
func NewControllerMedia() *ControllerMedia {
	return &ControllerMedia{}
}

func (c ControllerMedia) StreamMSE(w http.ResponseWriter, r *http.Request) {

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Error(err.Error())
		return
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf(err.Error())
		}
	}()

	entityId := r.PathValue("entity_id")
	token := r.URL.Query().Get("token")
	clientIp := ""
	channel := r.PathValue("channel")

	if !server.Storage.StreamChannelExist(entityId, channel) {
		//log.Error(media.ErrorStreamNotFound.Error())
		return
	}

	if !server.RemoteAuthorization("WS", entityId, channel, token, clientIp) {
		log.Error(server.ErrorStreamUnauthorized.Error())
		return
	}

	server.Storage.StreamChannelRun(entityId, channel)
	err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Error(err.Error())
		return
	}
	cid, ch, _, err := server.Storage.ClientAdd(entityId, channel, server.MSE)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer server.Storage.ClientDelete(entityId, cid, channel)
	codecs, err := server.Storage.StreamChannelCodecs(entityId, channel)
	if err != nil {
		log.Error(err.Error())
		return
	}
	muxerMSE := mp4f.NewMuxer(nil)
	err = muxerMSE.WriteHeader(codecs)
	if err != nil {
		log.Error(err.Error())
		return
	}
	meta, init := muxerMSE.GetInit(codecs)
	err = wsutil.WriteServerMessage(conn, ws.OpBinary, append([]byte{9}, meta...))
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = wsutil.WriteServerMessage(conn, ws.OpBinary, init)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var videoStart bool
	controlExit := make(chan struct{})
	noClient := time.NewTimer(10 * time.Second)
	go func() {
		var header ws.Header
		defer func() {
			close(controlExit)
		}()
		for {
			header, _, err = wsutil.NextReader(conn, ws.StateServerSide)
			if err != nil {
				return
			}
			switch header.OpCode {
			case ws.OpPong:
				noClient.Reset(10 * time.Second)
			case ws.OpClose:
				return
			}
		}
	}()
	noVideo := time.NewTimer(10 * time.Second)
	pingTicker := time.NewTicker(500 * time.Millisecond)
	defer pingTicker.Stop()
	var buf []byte
	for {
		select {

		case <-pingTicker.C:
			if err = conn.SetWriteDeadline(time.Now().Add(3 * time.Second)); err != nil {
				log.Error(err.Error())
				return
			}
			if buf, err = ws.CompileFrame(ws.NewPingFrame(nil)); err != nil {
				log.Error(err.Error())
				return
			}
			if _, err = conn.Write(buf); err != nil {
				log.Error(err.Error())
				return
			}
		case <-controlExit:
			return
		case <-noClient.C:
			return
		case <-noVideo.C:
			return
		case pck := <-ch:
			if pck.IsKeyFrame {
				noVideo.Reset(10 * time.Second)
				videoStart = true
			}
			if !videoStart {
				continue
			}
			var ready bool
			if ready, buf, err = muxerMSE.WritePacket(*pck, false); err != nil {
				log.Error(err.Error())
				return
			}
			if ready {
				if err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
					log.Error(err.Error())
					return
				}
				//err = websocket.Message.Send(ws, buf)
				if err = wsutil.WriteServerMessage(conn, ws.OpBinary, buf); err != nil {
					log.Error(err.Error())
					return
				}
			}
		}
	}
}

func (c ControllerMedia) StreamHLSLLInit(w http.ResponseWriter, r *http.Request) {
}

func (c ControllerMedia) StreamHLSLLM3U8(w http.ResponseWriter, r *http.Request) {
}

func (c ControllerMedia) StreamHLSLLM4Segment(w http.ResponseWriter, r *http.Request) {
}

func (c ControllerMedia) StreamHLSLLM4Fragment(w http.ResponseWriter, r *http.Request) {
}
