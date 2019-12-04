package models

// swagger:model
type FlowSubscription struct {
	Id    int64  `json:"id"`
	Topic string `json:"topic"`
}
