package core

type Telemetry interface {
	Broadcast(string)
	BroadcastOne(string, int64)
}