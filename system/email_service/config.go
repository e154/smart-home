package email_service

type EmailServiceConfig struct {
	Auth   string
	Pass   string
	Smtp   string
	Port   int
	Sender string
}

func NewEmailServiceConfig(auth, pass, smtp string, port int, send string) *EmailServiceConfig {
	return &EmailServiceConfig{
		Auth:   auth,
		Pass:   pass,
		Smtp:   smtp,
		Port:   port,
		Sender: send,
	}
}
