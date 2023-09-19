package controllers

import (
	"fmt"
	"github.com/deepch/vdk/format/mp4f"
	"github.com/e154/smart-home/system/media"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/labstack/echo/v4"
	"time"
)

// ControllerMedia ...
type ControllerMedia struct {
	*ControllerCommon
}

// NewControllerMedia ...
func NewControllerMedia(common *ControllerCommon) ControllerMedia {
	return ControllerMedia{
		ControllerCommon: common,
	}
}

func (c ControllerMedia) StreamMSE(ctx echo.Context) error {

	conn, _, _, err := ws.UpgradeHTTP(ctx.Request(), ctx.Response().Writer)
	if err != nil {
		return c.ERROR(ctx, err)
	}

	entityId := ctx.Param("entity_id")
	token := ctx.Param("token")
	clientIp := ctx.RealIP()
	channel := "0"

	fmt.Println(entityId, token, clientIp)

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error(err.Error())
		}
		log.Debug("Client Full Exit")
	}()
	if !media.Storage.StreamChannelExist(entityId, channel) {
		log.Error(media.ErrorStreamNotFound.Error())
		return c.ERROR(ctx, err)
	}

	if !media.RemoteAuthorization("WS", entityId, channel, token, clientIp) {
		log.Error(media.ErrorStreamUnauthorized.Error())
		return c.ERROR(ctx, err)
	}

	media.Storage.StreamChannelRun(entityId, channel)
	err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return c.ERROR(ctx, err)
	}
	cid, ch, _, err := media.Storage.ClientAdd(entityId, channel, media.MSE)
	if err != nil {
		return c.ERROR(ctx, err)
	}
	defer media.Storage.ClientDelete(entityId, cid, channel)
	codecs, err := media.Storage.StreamChannelCodecs(entityId, channel)
	if err != nil {
		return c.ERROR(ctx, err)
	}
	muxerMSE := mp4f.NewMuxer(nil)
	err = muxerMSE.WriteHeader(codecs)
	if err != nil {
		return c.ERROR(ctx, err)
	}
	meta, init := muxerMSE.GetInit(codecs)
	err = wsutil.WriteServerMessage(conn, ws.OpBinary, append([]byte{9}, meta...))
	if err != nil {
		return c.ERROR(ctx, err)
	}
	err = wsutil.WriteServerMessage(conn, ws.OpBinary, init)
	if err != nil {
		return c.ERROR(ctx, err)
	}
	var videoStart bool
	controlExit := make(chan bool, 10)
	noClient := time.NewTimer(10 * time.Second)
	go func() {
		defer func() {
			controlExit <- true
		}()
		for {
			header, _, err := wsutil.NextReader(conn, ws.StateServerSide)
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
	defer log.Debug("client exit")
	for {
		select {

		case <-pingTicker.C:
			err = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			if err != nil {
				return c.ERROR(ctx, err)
			}
			buf, err := ws.CompileFrame(ws.NewPingFrame(nil))
			if err != nil {
				return c.ERROR(ctx, err)
			}
			_, err = conn.Write(buf)
			if err != nil {
				return c.ERROR(ctx, err)
			}
		case <-controlExit:
			log.Info("Client Reader Exit")
			return c.ERROR(ctx, err)
		case <-noClient.C:
			log.Info("Client OffLine Exit")
			return c.ERROR(ctx, err)
		case <-noVideo.C:
			log.Info(media.ErrorStreamNoVideo.Error())
			return c.ERROR(ctx, err)
		case pck := <-ch:
			if pck.IsKeyFrame {
				noVideo.Reset(10 * time.Second)
				videoStart = true
			}
			if !videoStart {
				continue
			}
			ready, buf, err := muxerMSE.WritePacket(*pck, false)
			if err != nil {
				return c.ERROR(ctx, err)
			}
			if ready {
				err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					return c.ERROR(ctx, err)
				}
				//err = websocket.Message.Send(ws, buf)
				err = wsutil.WriteServerMessage(conn, ws.OpBinary, buf)
				if err != nil {
					return c.ERROR(ctx, err)
				}
			}
		}
	}
}
