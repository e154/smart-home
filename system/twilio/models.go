package twilio

type Balance struct {
	Currency   string `json:"currency"`
	Balance    string `json:"balance"`
	AccountSid string `json:"account_sid"`
}
