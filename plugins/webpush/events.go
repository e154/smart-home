package webpush

import m "github.com/e154/smart-home/models"

// EventAddWebPushSubscription ...
type EventAddWebPushSubscription struct {
	UserID       int64           `json:"user_id"`
	Subscription *m.Subscription `json:"subscription"`
}

// EventGetWebPushPublicKey ...
type EventGetWebPushPublicKey struct {
	UserID int64 `json:"user_id,omitempty"`
}

// EventNewWebPushPublicKey ...
type EventNewWebPushPublicKey struct {
	UserID    int64  `json:"user_id,omitempty"`
	PublicKey string `json:"public_key"`
}
