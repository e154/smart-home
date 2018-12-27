package models

import "time"

type NewDevice struct {
	Id          int64  `json:"id"`
	DeviceId    *int64 `json:"device_id"`
	NodeId      *int64  `json:"node_id"`
	Address     *int   `json:"address"`
	Baud        int    `json:"baud"`
	Sleep       int64  `json:"sleep"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	StopBite    int64  `json:"stop_bite"`
	Timeout     int64  `json:"timeout"`
	Tty         string `json:"tty"`
	IsGroup     bool   `json:"is_group"`
}

type UpdateDevice struct {
	Id          int64  `json:"id"`
	DeviceId    *int64 `json:"device_id"`
	NodeId      *int64  `json:"node_id"`
	Address     *int   `json:"address"`
	Baud        int    `json:"baud"`
	Sleep       int64  `json:"sleep"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	StopBite    int64  `json:"stop_bite"`
	Timeout     int64  `json:"timeout"`
	Tty         string `json:"tty"`
	IsGroup     bool   `json:"is_group"`
}

type Device struct {
	Id          int64     `json:"id"`
	DeviceId    *int64    `json:"device_id"`
	NodeId      *int64     `json:"node_id"`
	Address     *int      `json:"address"`
	Baud        int       `json:"baud"`
	Sleep       int64     `json:"sleep"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	StopBite    int64     `json:"stop_bite"`
	Timeout     int64     `json:"timeout"`
	Tty         string    `json:"tty"`
	IsGroup     bool      `json:"is_group"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Devices []*Device

type ResponseDevice struct {
	Code ResponseType `json:"code"`
	Data struct {
		Device *Device `json:"device"`
	} `json:"data"`
}

type DeviceListModel struct {
	Items []UserShotModel
	Meta  struct {
		Limit        int `json:"limit"`
		Offset       int `json:"offset"`
		ObjectsCount int `json:"objects_count"`
	}
}