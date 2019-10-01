package stream

type ConnectType string

const (
	WEBSOCK = ConnectType("websock")
)

const (
	Request       = "request"
	Response      = "response"
	StatusSuccess = "success"
	StatusError   = "error"
	Notify        = "notify"
	Broadcast     = "broadcast"
)

type BroadcastClient interface {
	Broadcast(message []byte)
}