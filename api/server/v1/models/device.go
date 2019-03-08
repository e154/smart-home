package models

import "time"

type NewDeviceNode struct {
	Id int64 `json:"id"`
}

type NewDeviceModel struct {
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


type DeviceModel struct {
	Id          int64            `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Node        *NodeModel       `json:"node"`
	Properties  map[string]interface{} `json:"properties"`
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

type ResponseDevice struct {
	Device *DeviceModel `json:"device"`
}

type DeviceListModel struct {
	Items []DeviceModel `json:"items"`
	Meta  struct {
		Limit        int `json:"limit"`
		Offset       int `json:"offset"`
		ObjectsCount int `json:"objects_count"`
	}
}

type SearchDeviceResponse struct {
	Devices []DeviceModel `json:"devices"`
}
