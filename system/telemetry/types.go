package telemetry

type ITelemetry interface {
	Broadcast(pack string)
	BroadcastOne(pack string, deviceId int64, elementName string)
}

