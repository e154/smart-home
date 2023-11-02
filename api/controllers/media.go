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

package controllers

import (
	"time"

	"github.com/deepch/vdk/format/mp4f"
	"github.com/e154/smart-home/system/media"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/labstack/echo/v4"
)

// ControllerMedia ...
type ControllerMedia struct {
	*ControllerCommon
}

// NewControllerMedia ...
func NewControllerMedia(common *ControllerCommon) *ControllerMedia {
	return &ControllerMedia{
		ControllerCommon: common,
	}
}

func (c ControllerMedia) StreamMSE(ctx echo.Context) error {

	conn, _, _, err := ws.UpgradeHTTP(ctx.Request(), ctx.Response().Writer)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf(err.Error())
		}
	}()

	entityId := ctx.Param("entity_id")
	token := ctx.Param("token")
	clientIp := ctx.RealIP()
	channel := ctx.Param("channel")

	if !media.Storage.StreamChannelExist(entityId, channel) {
		//log.Error(media.ErrorStreamNotFound.Error())
		return nil
	}

	if !media.RemoteAuthorization("WS", entityId, channel, token, clientIp) {
		log.Error(media.ErrorStreamUnauthorized.Error())
		return nil
	}

	media.Storage.StreamChannelRun(entityId, channel)
	err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	cid, ch, _, err := media.Storage.ClientAdd(entityId, channel, media.MSE)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	defer media.Storage.ClientDelete(entityId, cid, channel)
	codecs, err := media.Storage.StreamChannelCodecs(entityId, channel)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	muxerMSE := mp4f.NewMuxer(nil)
	err = muxerMSE.WriteHeader(codecs)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	meta, init := muxerMSE.GetInit(codecs)
	err = wsutil.WriteServerMessage(conn, ws.OpBinary, append([]byte{9}, meta...))
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	err = wsutil.WriteServerMessage(conn, ws.OpBinary, init)
	if err != nil {
		log.Error(err.Error())
		return nil
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
				return nil
			}
			if buf, err = ws.CompileFrame(ws.NewPingFrame(nil)); err != nil {
				log.Error(err.Error())
				return nil
			}
			if _, err = conn.Write(buf); err != nil {
				log.Error(err.Error())
				return nil
			}
		case <-controlExit:
			return nil
		case <-noClient.C:
			return nil
		case <-noVideo.C:
			return nil
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
				return nil
			}
			if ready {
				if err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
					log.Error(err.Error())
					return nil
				}
				//err = websocket.Message.Send(ws, buf)
				if err = wsutil.WriteServerMessage(conn, ws.OpBinary, buf); err != nil {
					log.Error(err.Error())
					return nil
				}
			}
		}
	}
}

func (c ControllerMedia) StreamHLSLLInit(ctx echo.Context) error {
	return nil
}

func (c ControllerMedia) StreamHLSLLM3U8(ctx echo.Context) error {
	return nil
}

func (c ControllerMedia) StreamHLSLLM4Segment(ctx echo.Context) error {
	return nil
}

func (c ControllerMedia) StreamHLSLLM4Fragment(ctx echo.Context) error {
	return nil
}
