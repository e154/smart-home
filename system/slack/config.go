package slack

type SlackConfig struct {
	Token    string
	UserName string
}

func NewSlackConfig(slackToken, userName string) *SlackConfig {
	return &SlackConfig{
		Token:    slackToken,
		UserName: userName,
	}
}
