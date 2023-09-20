package media

import (
	"github.com/e154/smart-home/common/apperr"
	"net"
	"sync"
	"time"

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
	Success                         = "success"
	ErrorStreamNotFound             = apperr.New("stream not found", apperr.ErrNotFound)
	ErrorStreamAlreadyExists        = apperr.New("stream already exists", apperr.ErrInvalidRequest)
	ErrorStreamChannelAlreadyExists = apperr.New("stream channel already exists", apperr.ErrInvalidRequest)
	ErrorStreamNotHLSSegments       = apperr.New("stream hls not ts seq found", apperr.ErrNotFound)
	ErrorStreamNoVideo              = apperr.New("stream no video", apperr.ErrNotFound)
	ErrorStreamNoClients            = apperr.New("stream no clients", apperr.ErrNotFound)
	ErrorStreamRestart              = apperr.New("stream restart", apperr.ErrInternal)
	ErrorStreamStopCoreSignal       = apperr.New("stream stop core signal", apperr.ErrInternal)
	ErrorStreamStopRTSPSignal       = apperr.New("stream stop rtsp signal", apperr.ErrInternal)
	ErrorStreamChannelNotFound      = apperr.New("stream channel not found", apperr.ErrNotFound)
	ErrorStreamChannelCodecNotFound = apperr.New("stream channel codec not ready, possible stream offline", apperr.ErrInternal)
	ErrorStreamsLen0                = apperr.New("streams len zero", apperr.ErrInternal)
	ErrorStreamUnauthorized         = apperr.New("stream request unauthorized", apperr.ErrUnauthorized)
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
