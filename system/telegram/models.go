package telegram

type Command struct {
	UserName, Text string
	ChatId         int64
}
