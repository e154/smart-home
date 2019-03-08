package models

import "time"

type NewDeviceNode struct {
	Id int64 `json:"id"`
}

type NewDevice struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	Device      *ParentDevice          `json:"device"`
	Type        string                 `json:"type"`
	Node        *NewDeviceNode         `json:"node"`
	Properties  map[string]interface{} `json:"properties"`
}

type UpdateDevice struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	Device      *ParentDevice          `json:"device"`
	Type        string                 `json:"type"`
	Node        *NewDeviceNode         `json:"node"`
	Properties  map[string]interface{} `json:"properties"`
}

type ParentDevice struct {
	Id int64 `json:"id"`
}

type DeviceProperties map[string]interface{}

type Device struct {
	Id          int64            `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Node        *NodeModel       `json:"node"`
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

type Devices []*Device

type ResponseDevice struct {
	Code ResponseType `json:"code"`
	Data struct {
		Device *Device `json:"device"`
	} `json:"data"`
}

type DeviceListModel struct {
	Items []Device
	Meta  struct {
		Limit        int `json:"limit"`
		Offset       int `json:"offset"`
		ObjectsCount int `json:"objects_count"`
	}
}

type SearchDeviceResponse struct {
	Devices []Device `json:"devices"`
}
