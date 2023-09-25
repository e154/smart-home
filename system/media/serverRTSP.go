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
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	Version   = "RTSP/1.0"
	UserAgent = "Lavf58.29.100"
	Session   = "000a959d6816"
)
var (
	OPTIONS  = "OPTIONS"
	DESCRIBE = "DESCRIBE"
	SETUP    = "SETUP"
	PLAY     = "PLAY"
	TEARDOWN = "TEARDOWN"
)

// RTSP response status codes
const (
	StatusContinue                      = 100
	StatusOK                            = 200
	StatusCreated                       = 201
	StatusLowOnStorageSpace             = 250
	StatusMultipleChoices               = 300
	StatusMovedPermanently              = 301
	StatusMovedTemporarily              = 302
	StatusSeeOther                      = 303
	StatusNotModified                   = 304
	StatusUseProxy                      = 305
	StatusBadRequest                    = 400
	StatusUnauthorized                  = 401
	StatusPaymentRequired               = 402
	StatusForbidden                     = 403
	StatusNotFound                      = 404
	StatusMethodNotAllowed              = 405
	StatusNotAcceptable                 = 406
	StatusProxyAuthenticationRequired   = 407
	StatusRequestTimeout                = 408
	StatusGone                          = 410
	StatusLengthRequired                = 411
	StatusPreconditionFailed            = 412
	StatusRequestEntityTooLarge         = 413
	StatusRequestURITooLong             = 414
	StatusUnsupportedMediaType          = 415
	StatusInvalidparameter              = 451
	StatusIllegalConferenceIdentifier   = 452
	StatusNotEnoughBandwidth            = 453
	StatusSessionNotFound               = 454
	StatusMethodNotValidInThisState     = 455
	StatusHeaderFieldNotValid           = 456
	StatusInvalidRange                  = 457
	StatusParameterIsReadOnly           = 458
	StatusAggregateOperationNotAllowed  = 459
	StatusOnlyAggregateOperationAllowed = 460
	StatusUnsupportedTransport          = 461
	StatusDestinationUnreachable        = 462
	StatusInternalServerError           = 500
	StatusNotImplemented                = 501
	StatusBadGateway                    = 502
	StatusServiceUnavailable            = 503
	StatusGatewayTimeout                = 504
	StatusRTSPVersionNotSupported       = 505
	StatusOptionNotsupport              = 551
)

func StatusText(code int) string {
	return statusText[code]
}

var statusText = map[int]string{
	StatusContinue:                      "Continue",
	StatusOK:                            "OK",
	StatusCreated:                       "Created",
	StatusLowOnStorageSpace:             "Low on Storage Space",
	StatusMultipleChoices:               "Multiple Choices",
	StatusMovedPermanently:              "Moved Permanently",
	StatusMovedTemporarily:              "Moved Temporarily",
	StatusSeeOther:                      "See Other",
	StatusNotModified:                   "Not Modified",
	StatusUseProxy:                      "Use Proxy",
	StatusBadRequest:                    "Bad Request",
	StatusUnauthorized:                  "Unauthorized",
	StatusPaymentRequired:               "Payment Required",
	StatusForbidden:                     "Forbidden",
	StatusNotFound:                      "Not Found",
	StatusMethodNotAllowed:              "Method Not Allowed",
	StatusNotAcceptable:                 "Not Acceptable",
	StatusProxyAuthenticationRequired:   "Proxy Authentication Required",
	StatusRequestTimeout:                "Request Time-out",
	StatusGone:                          "Gone",
	StatusLengthRequired:                "Length Required",
	StatusPreconditionFailed:            "Precondition Failed",
	StatusRequestEntityTooLarge:         "Request Entity Too Large",
	StatusRequestURITooLong:             "Request-URI Too Large",
	StatusUnsupportedMediaType:          "Unsupported Media Type",
	StatusInvalidparameter:              "Parameter Not Understood",
	StatusIllegalConferenceIdentifier:   "Conference Not Found",
	StatusNotEnoughBandwidth:            "Not Enough Bandwidth",
	StatusSessionNotFound:               "Session Not Found",
	StatusMethodNotValidInThisState:     "Method Not Valid in This State",
	StatusHeaderFieldNotValid:           "Header Field Not Valid for Resource",
	StatusInvalidRange:                  "Invalid Range",
	StatusParameterIsReadOnly:           "Parameter Is Read-Only",
	StatusAggregateOperationNotAllowed:  "Aggregate operation not allowed",
	StatusOnlyAggregateOperationAllowed: "Only aggregate operation allowed",
	StatusUnsupportedTransport:          "Unsupported transport",
	StatusDestinationUnreachable:        "Destination unreachable",
	StatusInternalServerError:           "Internal Server Error",
	StatusNotImplemented:                "Not Implemented",
	StatusBadGateway:                    "Bad Gateway",
	StatusServiceUnavailable:            "Service Unavailable",
	StatusGatewayTimeout:                "Gateway Time-out",
	StatusRTSPVersionNotSupported:       "RTSP Version not supported",
	StatusOptionNotsupport:              "Option not supported",
}

// RTSPServer func
func RTSPServer() {
	log.Info("Server RTSP start")
	l, err := net.Listen("tcp", Storage.ServerRTSPPort())
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer func() {
		err := l.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error(err.Error())
			return
		}
		go RTSPServerClientHandle(conn)
	}
}

// RTSPServerClientHandle func
func RTSPServerClientHandle(conn net.Conn) {
	buf := make([]byte, 4096)
	token, uuid, channel, in, cSEQ := "", "", "0", 0, 0
	var playStarted bool
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(err.Error())
		}

	}()
	err := conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Error(err.Error())
		return
	}
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Error(err.Error())
			return
		}
		cSEQ = parsecSEQ(buf[:n])
		stage, err := parseStage(buf[:n])
		if err != nil {
			log.Error(err.Error())
		}
		err = conn.SetDeadline(time.Now().Add(60 * time.Second))
		log.Debug(string(buf[:n]))

		if err != nil {
			log.Error(err.Error())
			return
		}

		switch stage {
		case OPTIONS:
			if playStarted {
				err = RTSPServerClientResponse(uuid, channel, conn, 200, map[string]string{"CSeq": strconv.Itoa(cSEQ), "Public": "DESCRIBE, SETUP, TEARDOWN, PLAY"})
				if err != nil {
					return
				}
				continue
			}
			uuid, channel, token, err = parseStreamChannel(buf[:n])
			if err != nil {
				log.Error(err.Error())
				return
			}
			if !Storage.StreamChannelExist(uuid, channel) {
				log.Error(ErrorStreamNotFound.Error())
				err = RTSPServerClientResponse(uuid, channel, conn, 404, map[string]string{"CSeq": strconv.Itoa(cSEQ)})
				if err != nil {
					return
				}
				return
			}

			if !RemoteAuthorization("RTSP", uuid, channel, token, conn.RemoteAddr().String()) {
				log.Error(ErrorStreamUnauthorized.Error())
				err = RTSPServerClientResponse(uuid, channel, conn, 401, map[string]string{"CSeq": strconv.Itoa(cSEQ)})
				if err != nil {
					return
				}
				return
			}

			Storage.StreamChannelRun(uuid, channel)
			err = RTSPServerClientResponse(uuid, channel, conn, 200, map[string]string{"CSeq": strconv.Itoa(cSEQ), "Public": "DESCRIBE, SETUP, TEARDOWN, PLAY"})
			if err != nil {
				return
			}
		case SETUP:
			if !strings.Contains(string(buf[:n]), "interleaved") {
				err = RTSPServerClientResponse(uuid, channel, conn, 461, map[string]string{"CSeq": strconv.Itoa(cSEQ)})
				if err != nil {
					return
				}
				continue
			}
			err = RTSPServerClientResponse(uuid, channel, conn, 200, map[string]string{"CSeq": strconv.Itoa(cSEQ), "User-Agent:": UserAgent, "Session": Session, "Transport": "RTP/AVP/TCP;unicast;interleaved=" + strconv.Itoa(in) + "-" + strconv.Itoa(in+1)})
			if err != nil {
				return
			}
			in = in + 2
		case DESCRIBE:
			sdp, err := Storage.StreamChannelSDP(uuid, channel)
			if err != nil {
				log.Error(err.Error())
				return
			}
			err = RTSPServerClientResponse(uuid, channel, conn, 200, map[string]string{"CSeq": strconv.Itoa(cSEQ), "User-Agent:": UserAgent, "Session": Session, "Content-Type": "application/sdp\r\nContent-Length: " + strconv.Itoa(len(sdp)), "sdp": string(sdp)})
			if err != nil {
				return
			}
		case PLAY:
			err = RTSPServerClientResponse(uuid, channel, conn, 200, map[string]string{"CSeq": strconv.Itoa(cSEQ), "User-Agent:": UserAgent, "Session": Session})
			if err != nil {
				return
			}
			playStarted = true
			go RTSPServerClientPlay(uuid, channel, conn)
		case TEARDOWN:
			err = RTSPServerClientResponse(uuid, channel, conn, 200, map[string]string{"CSeq": strconv.Itoa(cSEQ), "User-Agent:": UserAgent, "Session": Session})
			if err != nil {
				return
			}
			return
		default:
			log.Debugf("stage bad %v", stage)
		}
	}
}

// handleRTSPServerPlay func
func RTSPServerClientPlay(uuid string, channel string, conn net.Conn) {
	cid, _, ch, err := Storage.ClientAdd(uuid, channel, RTSP)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer func() {
		Storage.ClientDelete(uuid, cid, channel)
		log.Info("Client offline")
		err := conn.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	noVideo := time.NewTimer(10 * time.Second)

	for {
		select {
		case <-noVideo.C:
			return
		case pck := <-ch:
			noVideo.Reset(10 * time.Second)
			_, err := conn.Write(*pck)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	}
}

// handleRTSPServerPlay func
func RTSPServerClientResponse(uuid string, channel string, conn net.Conn, status int, headers map[string]string) error {
	var sdp string
	builder := bytes.Buffer{}
	builder.WriteString(fmt.Sprintf(Version+" %d %s\r\n", status, StatusText(status)))
	for k, v := range headers {
		if k == "sdp" {
			sdp = v
			continue
		}
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	builder.WriteString(fmt.Sprintf("\r\n"))
	builder.WriteString(sdp)
	log.Debug(builder.String())
	if _, err := conn.Write(builder.Bytes()); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// parsecSEQ func
func parsecSEQ(buf []byte) int {
	return stringToInt(stringInBetween(string(buf), "CSeq: ", "\r\n"))
}

// parseStage func
func parseStage(buf []byte) (string, error) {
	st := strings.Split(string(buf), " ")
	if len(st) > 0 {
		return st[0], nil
	}
	return "", errors.New("parse stage error " + string(buf))
}

// parseStreamChannel func
func parseStreamChannel(buf []byte) (string, string, string, error) {

	var token string

	uri := stringInBetween(string(buf), " ", " ")
	u, err := url.Parse(uri)
	if err == nil {
		token = u.Query().Get("token")
		uri = u.Path
	}

	st := strings.Split(uri, "/")

	if len(st) >= 3 {
		return st[1], st[2], token, nil
	}

	return "", "0", token, errors.New("parse stream error " + string(buf))
}
