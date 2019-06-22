package gate_client

type IWsCallback interface {
	onMessage(payload []byte)
	onConnected()
	onClosed()
}

const (
	ClientTypeServer = "server"
)