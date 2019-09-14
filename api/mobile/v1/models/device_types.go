package models

type DevModBusConfig struct {
	SlaveId  int    `json:"slave_id" mapstructure:"slave_id"`   // 1-32
	Baud     int    `json:"baud"`                               // 9600, 19200, ...
	DataBits int    `json:"data_bits" mapstructure:"data_bits"` // 5-9
	StopBits int    `json:"stop_bits" mapstructure:"stop_bits"` // 1, 2
	Parity   string `json:"parity"`                             // none, odd, even
	Timeout  int    `json:"timeout"`                            // milliseconds
}

type DevSmartBusConfig struct {
	Baud     int `json:"baud" valid:"Required"`
	Device   int `json:"device"`
	Timeout  int `json:"timeout" valid:"Required"`
	StopBits int `json:"stop_bits" valid:"Required" mapstructure:"stop_bits"`
	Sleep    int `json:"sleep"`
}

type DevCommandConfig struct {
}

// An AllOfModel is composed out of embedded structs but it should build
// an allOf property
type DeviceProperties struct {
	// swagger:allOf
	DevModBusConfig
	// swagger:allOf
	DevSmartBusConfig
	// swagger:allOf
	DevCommandConfig
}
