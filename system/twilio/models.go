package twilio

type Balance struct {
	Currency   string `json:"currency"`
	Balance    string `json:"balance"`
	AccountSid string `json:"account_sid"`
}

const (
	StatusAccepted    = "accepted"
	StatusQueued      = "queued"
	StatusSending     = "sending"
	StatusReceiving   = "receiving"
	StatusReceived    = "received"
	StatusDelivered   = "delivered"
	StatusUndelivered = "undelivered"
	StatusSent        = "sent"
	StatusFailed      = "failed"
)
