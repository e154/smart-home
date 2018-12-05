package core

import (
	"encoding/json"
	"github.com/e154/smart-home/common"
	"time"
)

type NodeMessage struct {
	DeviceId   int64             `json:"device_id"`
	DeviceType common.DeviceType `json:"device_type"`
	Properties json.RawMessage   `json:"properties"`
	Command    json.RawMessage   `json:"command"`
}

type NodeResponse struct {
	DeviceId   int64             `json:"device_id"`
	DeviceType common.DeviceType `json:"device_type"`
	Properties json.RawMessage   `json:"properties"`
	Response   json.RawMessage   `json:"response"`
	Status     string            `json:"status"`
	Time       float64
}

type NodeStatus string

type NodeStatModel struct {
	Status    NodeStatus `json:"status"`
	Thread    int        `json:"thread"`
	Rps       int64      `json:"rps"`
	Min       int64      `json:"min"`
	Max       int64      `json:"max"`
	StartedAt time.Time  `json:"started_at"`
}
