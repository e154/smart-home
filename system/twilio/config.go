package twilio

type TWConfig struct {
	from      string
	sid       string
	authToken string
}

func NewTWConfig(from, sid, authToken string) *TWConfig {
	return &TWConfig{
		from:      from,
		sid:       sid,
		authToken: authToken,
	}
}
