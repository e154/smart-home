package events

type EventServiceStarted struct {
	Service string `json:"service"`
}

type EventServiceStopped struct {
	Service string `json:"service"`
}

type EventServiceRestarted struct {
	Service string `json:"service"`
}
