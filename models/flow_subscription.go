package models

type FlowSubscription struct {
	Id     int64  `json:"id"`
	FlowId int64  `json:"flow_id"`
	Topic  string `json:"topic"`
}
