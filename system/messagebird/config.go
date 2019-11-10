package messagebird

type MBClientConfig struct {
	AccessKey  string
	Name       string
	CanSendSms bool
}

func NewMBClientConfig(accessKey, name string) *MBClientConfig {
	return &MBClientConfig{
		AccessKey: accessKey,
		Name:      name,
	}
}
