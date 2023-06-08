package models

import (
	"time"
)

// Keys are the base64 encoded values from PushSubscription.getKey()
type Keys struct {
	Auth   string `json:"auth"`
	P256dh string `json:"p256dh"`
}

// Subscription represents a PushSubscription object from the Push API
type Subscription struct {
	Endpoint string `json:"endpoint"`
	Keys     Keys   `json:"keys"`
}

// UserDevice ...
type UserDevice struct {
	Id           int64
	UserId       int64
	Subscription *Subscription
	CreatedAt    time.Time
}
