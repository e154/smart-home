package events

// EventUpdateUserLocation ...
type EventUpdateUserLocation struct {
	UserID   int64   `json:"user_id"`
	Lat      float32 `json:"lat"`
	Lon      float32 `json:"lon"`
	Accuracy float32 `json:"accuracy"`
}

type EventDirectMessage struct {
	UserID  int64       `json:"user_id"`
	Query   string      `json:"query"`
	Message interface{} `json:"message"`
}
