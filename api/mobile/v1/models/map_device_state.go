package models

type MapDeviceState struct {
	Id            int64        `json:"id"`
	DeviceState   *DeviceState `json:"device_state"`
	DeviceStateId int64        `json:"device_state_id" valid:"Required"`
	MapDeviceId   int64        `json:"map_device_id" valid:"Required"`
	Image         *Image       `json:"image"`
	ImageId       int64        `json:"image_id" valid:"Required"`
	Style         string       `json:"style"`
}
