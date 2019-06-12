package models

import "time"

type NewDeviceNode struct {
	Id int64 `json:"id"`
}

// swagger:model
type NewDevice struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Device      *ParentDevice    `json:"device"`
	Type        string           `json:"type"`
	Node        *NewDeviceNode   `json:"node"`
	Properties  DeviceProperties `json:"properties"`
}

// swagger:model
type UpdateDevice struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Device      *ParentDevice    `json:"device"`
	Type        string           `json:"type"`
	Node        *NewDeviceNode   `json:"node"`
	Properties  DeviceProperties `json:"properties"`
}

type ParentDevice struct {
	Id int64 `json:"id"`
}

// swagger:model
type Device struct {
	Id          int64            `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Node        *Node            `json:"node"`
	Properties  DeviceProperties `json:"properties"`
	Type        string           `json:"type"`
	Status      string           `json:"status"`
	IsGroup     bool             `json:"is_group"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Actions     []DeviceAction   `json:"actions"`
	States      []DeviceState    `json:"states"`
	Device      *ParentDevice    `json:"device"`
	DeviceId    *int64           `json:"device_id"`
}

// swagger:model
type DeviceShort struct {
	Id          int64            `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Type        string           `json:"type"`
}
