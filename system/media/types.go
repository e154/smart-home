package media

type EventUpdateList struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
}

type EventRemoveList struct {
	Name string `json:"name"`
}
