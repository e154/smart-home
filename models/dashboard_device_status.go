package models

type DashboardDeviceStatus struct {
	Id          int64  `json:"id"`
	Description string `json:"description"`
	SystemName  string `json:"system_name"`
	DeviceId    int64  `json:"device_id"`
}
