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
	"net"
	"sync"
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/deepch/vdk/av"
)

var Storage = NewStreamCore()

// Default stream  type
const (
	MSE = iota
	WEBRTC
	RTSP
)

// Default stream status type
const (
	OFFLINE = iota
	ONLINE
)

// Default stream errors
var (
	ErrorStreamNotFound             = apperr.ErrorWithCode("STREAM_NOT_FOUND", "stream not found", apperr.ErrNotFound)
	ErrorStreamAlreadyExists        = apperr.ErrorWithCode("STREAM_ALREADY_EXISTS", "stream already exists", apperr.ErrInvalidRequest)
	ErrorStreamChannelAlreadyExists = apperr.ErrorWithCode("STREAM_CHANNEL_ALREADY_EXISTS", "stream channel already exists", apperr.ErrInvalidRequest)
	ErrorStreamNotHLSSegments       = apperr.ErrorWithCode("STREAM_HLS_NOT_TS_SEQ_FOUND", "stream hls not ts seq found", apperr.ErrNotFound)
	ErrorStreamNoVideo              = apperr.ErrorWithCode("STREAM_NO_VIDEO", "stream no video", apperr.ErrNotFound)
	ErrorStreamNoClients            = apperr.ErrorWithCode("STREAM_NO_CLIENTS", "stream no clients", apperr.ErrNotFound)
	ErrorStreamRestart              = apperr.ErrorWithCode("STREAM_RESTART", "stream restart", apperr.ErrInternal)
	ErrorStreamStopCoreSignal       = apperr.ErrorWithCode("STREAM_STOP_CORE_SIGNAL", "stream stop core signal", apperr.ErrInternal)
	ErrorStreamStopRTSPSignal       = apperr.ErrorWithCode("STREAM_STOP_RTSP_SIGNAL", "stream stop rtsp signal", apperr.ErrInternal)
	ErrorStreamChannelNotFound      = apperr.ErrorWithCode("STREAM_CHANNEL_NOT_FOUND", "stream channel not found", apperr.ErrNotFound)
	ErrorStreamChannelCodecNotFound = apperr.ErrorWithCode("STREAM_CHANNEL_CODEC_NOT_READY,_POSSIBLE_STREAM_OFFLINE", "stream channel codec not ready, possible stream offline", apperr.ErrInternal)
	ErrorStreamsLen0                = apperr.ErrorWithCode("STREAMS_LEN_ZERO", "streams len zero", apperr.ErrInternal)
	ErrorStreamUnauthorized         = apperr.ErrorWithCode("STREAM_REQUEST_UNAUTHORIZED", "stream request unauthorized", apperr.ErrUnauthorized)
)

// StorageST main storage struct
type StorageST struct {
	mutex           sync.RWMutex
	Server          ServerST            `json:"server" groups:"api,config"`
	Streams         map[string]StreamST `json:"streams,omitempty" groups:"api,config"`
	ChannelDefaults ChannelST           `json:"channel_defaults,omitempty" groups:"api,config"`
}

// ServerST server storage section
type ServerST struct {
	Debug bool `json:"debug" groups:"api,config"`
	//LogLevel           logrus.Level `json:"log_level" groups:"api,config"`
	HTTPDemo           bool     `json:"http_demo" groups:"api,config"`
	HTTPDebug          bool     `json:"http_debug" groups:"api,config"`
	HTTPLogin          string   `json:"http_login" groups:"api,config"`
	HTTPPassword       string   `json:"http_password" groups:"api,config"`
	HTTPDir            string   `json:"http_dir" groups:"api,config"`
	HTTPPort           string   `json:"http_port" groups:"api,config"`
	RTSPPort           string   `json:"rtsp_port" groups:"api,config"`
	HTTPS              bool     `json:"https" groups:"api,config"`
	HTTPSPort          string   `json:"https_port" groups:"api,config"`
	HTTPSCert          string   `json:"https_cert" groups:"api,config"`
	HTTPSKey           string   `json:"https_key" groups:"api,config"`
	HTTPSAutoTLSEnable bool     `json:"https_auto_tls" groups:"api,config"`
	HTTPSAutoTLSName   string   `json:"https_auto_tls_name" groups:"api,config"`
	ICEServers         []string `json:"ice_servers" groups:"api,config"`
	ICEUsername        string   `json:"ice_username" groups:"api,config"`
	ICECredential      string   `json:"ice_credential" groups:"api,config"`
	Token              Token    `json:"token,omitempty" groups:"api,config"`
	WebRTCPortMin      uint16   `json:"webrtc_port_min" groups:"api,config"`
	WebRTCPortMax      uint16   `json:"webrtc_port_max" groups:"api,config"`
}

// Token auth
type Token struct {
	Enable  bool   `json:"enable" groups:"api,config"`
	Backend string `json:"backend" groups:"api,config"`
}

// ServerST stream storage section
type StreamST struct {
	Name     string               `json:"name,omitempty" groups:"api,config"`
	Channels map[string]ChannelST `json:"channels,omitempty" groups:"api,config"`
}

type ChannelST struct {
	Name               string `json:"name,omitempty" groups:"api,config"`
	URL                string `json:"url,omitempty" groups:"api,config"`
	OnDemand           bool   `json:"on_demand,omitempty" groups:"api,config"`
	Debug              bool   `json:"debug,omitempty" groups:"api,config"`
	Status             int    `json:"status,omitempty" groups:"api"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify,omitempty" groups:"api,config"`
	Audio              bool   `json:"audio,omitempty" groups:"api,config"`
	runLock            bool
	codecs             []av.CodecData
	sdp                []byte
	signals            chan int
	hlsSegmentBuffer   map[int]SegmentOld
	hlsSegmentNumber   int
	clients            map[string]ClientST
	ack                time.Time
	hlsMuxer           *MuxerHLS `json:"-"`
}

// ClientST client storage section
type ClientST struct {
	mode              int
	signals           chan int
	outgoingAVPacket  chan *av.Packet
	outgoingRTPPacket chan *[]byte
	socket            net.Conn
}

// SegmentOld HLS cache section
type SegmentOld struct {
	dur  time.Duration
	data []*av.Packet
}
