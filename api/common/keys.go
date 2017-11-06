package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func ComputeHmac256() string {
	var message string = "token"
	var secret string = RandomString(255)

	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))

	return hex.EncodeToString(h.Sum(nil))
}