package mqtt_client

type Config struct {
	KeepAlive      int    `json:"keep_alive"`
	PingTimeout    int    `json:"ping_timeout"`
	Broker         string `json:"broker"`
	ClientID       string `json:"client_id"`
	ConnectTimeout int    `json:"connect_timeout"`
	CleanSession   bool   `json:"clean_session"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Qos            byte   `json:"qos"`
}
