// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package webpush

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"math/big"
	"strconv"
	"strings"
	"time"

	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/web"

	"golang.org/x/crypto/hkdf"
)

const MaxRecordSize uint32 = 4096

var ErrMaxPadExceeded = errors.New("payload has exceeded the maximum length")

// saltFunc generates a salt of 16 bytes
var saltFunc = func() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return salt, err
	}

	return salt, nil
}

// Options are config and extra params needed to send a notification
type Options struct {
	Crawler         web.Crawler // Will replace with *http.Client by default if not included
	RecordSize      uint32      // Limit the record size
	Subscriber      string      // Sub in VAPID JWT token
	Topic           string      // Set the Topic header to collapse a pending messages (Optional)
	TTL             int         // Set the TTL on the endpoint POST request
	Urgency         Urgency     // Set the Urgency header to change a message priority (Optional)
	VAPIDPublicKey  string      // VAPID public key, passed in VAPID Authorization header
	VAPIDPrivateKey string      // VAPID private key, used to sign VAPID JWT token
}

// SendNotification calls SendNotificationWithContext with default context for backwards-compatibility
func SendNotification(message []byte, s *m.Subscription, options *Options) (int, []byte, error) {
	return SendNotificationWithContext(context.Background(), message, s, options)
}

// SendNotificationWithContext sends a push notification to a subscription's endpoint
// Message Encryption for Web Push, and VAPID protocols.
// FOR MORE INFORMATION SEE RFC8291: https://datatracker.ietf.org/doc/rfc8291
func SendNotificationWithContext(ctx context.Context, message []byte, s *m.Subscription, options *Options) (code int, body []byte, err error) {
	// Authentication secret (auth_secret)
	var authSecret []byte
	authSecret, err = decodeSubscriptionKey(s.Keys.Auth)
	if err != nil {
		return
	}

	// dh (Diffie Hellman)
	var dh []byte
	if dh, err = decodeSubscriptionKey(s.Keys.P256dh); err != nil {
		return
	}

	// Generate 16 byte salt
	var salt []byte
	if salt, err = saltFunc(); err != nil {
		return
	}

	// Create the ecdh_secret shared key pair
	curve := elliptic.P256()

	// Application server key pairs (single use)
	var localPrivateKey []byte
	var x, y *big.Int
	localPrivateKey, x, y, err = elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		return
	}

	localPublicKey := elliptic.Marshal(curve, x, y)

	// Combine application keys with dh
	sharedX, sharedY := elliptic.Unmarshal(curve, dh)
	if sharedX == nil {
		err = errors.New("Unmarshal Error: Public key is not a valid point on the curve")
		return
	}

	sx, _ := curve.ScalarMult(sharedX, sharedY, localPrivateKey)
	sharedECDHSecret := sx.Bytes()

	hash := sha256.New

	// ikm
	prkInfoBuf := bytes.NewBuffer([]byte("WebPush: info\x00"))
	prkInfoBuf.Write(dh)
	prkInfoBuf.Write(localPublicKey)

	prkHKDF := hkdf.New(hash, sharedECDHSecret, authSecret, prkInfoBuf.Bytes())
	var ikm []byte
	if ikm, err = getHKDFKey(prkHKDF, 32); err != nil {
		return
	}

	// Derive Content Encryption Key
	contentEncryptionKeyInfo := []byte("Content-Encoding: aes128gcm\x00")
	contentHKDF := hkdf.New(hash, ikm, salt, contentEncryptionKeyInfo)
	var contentEncryptionKey []byte
	if contentEncryptionKey, err = getHKDFKey(contentHKDF, 16); err != nil {
		return
	}

	// Derive the Nonce
	nonceInfo := []byte("Content-Encoding: nonce\x00")
	nonceHKDF := hkdf.New(hash, ikm, salt, nonceInfo)
	var nonce []byte
	if nonce, err = getHKDFKey(nonceHKDF, 12); err != nil {
		return
	}

	// Cipher
	var c cipher.Block
	if c, err = aes.NewCipher(contentEncryptionKey); err != nil {
		return
	}

	var gcm cipher.AEAD
	if gcm, err = cipher.NewGCM(c); err != nil {
		return
	}

	// Get the record size
	recordSize := options.RecordSize
	if recordSize == 0 {
		recordSize = MaxRecordSize
	}

	recordLength := int(recordSize) - 16

	// Encryption Content-Coding Header
	recordBuf := bytes.NewBuffer(salt)

	rs := make([]byte, 4)
	binary.BigEndian.PutUint32(rs, recordSize)

	recordBuf.Write(rs)
	recordBuf.Write([]byte{byte(len(localPublicKey))})
	recordBuf.Write(localPublicKey)

	// Data
	dataBuf := bytes.NewBuffer(message)

	// Pad content to max record size - 16 - header
	// Padding ending delimeter
	dataBuf.Write([]byte("\x02"))
	if err = pad(dataBuf, recordLength-recordBuf.Len()); err != nil {
		return
	}

	// Compose the ciphertext
	ciphertext := gcm.Seal([]byte{}, nonce, dataBuf.Bytes(), nil)
	recordBuf.Write(ciphertext)

	req := web.Request{
		Method: "POST",
		Url:    s.Endpoint,
		Body:   recordBuf.Bytes(),
		Headers: []map[string]string{
			{"Content-Encoding": "aes128gcm"},
			{"Content-Length": strconv.Itoa(len(ciphertext))},
			{"Content-Type": "application/octet-stream"},
			{"TTL": strconv.Itoa(options.TTL)},
		},
		Timeout: time.Second * 5,
	}

	// Сheck the optional headers
	if len(options.Topic) > 0 {
		req.Headers = append(req.Headers, map[string]string{"Topic": options.Topic})
	}
	if isValidUrgency(options.Urgency) {
		req.Headers = append(req.Headers, map[string]string{"Urgency": string(options.Urgency)})
	}
	// Get VAPID Authorization header
	var vapidAuthHeader string
	vapidAuthHeader, err = getVAPIDAuthorizationHeader(
		s.Endpoint,
		options.Subscriber,
		options.VAPIDPublicKey,
		options.VAPIDPrivateKey,
	)
	if err != nil {
		return
	}

	req.Headers = append(req.Headers, map[string]string{"Authorization": vapidAuthHeader})

	code, body, err = options.Crawler.Probe(req)

	return
}

// decodeSubscriptionKey decodes a base64 subscription key.
// if necessary, add "=" padding to the key for URL decode
func decodeSubscriptionKey(key string) ([]byte, error) {
	// "=" padding
	buf := bytes.NewBufferString(key)
	if rem := len(key) % 4; rem != 0 {
		buf.WriteString(strings.Repeat("=", 4-rem))
	}

	bytes, err := base64.StdEncoding.DecodeString(buf.String())
	if err == nil {
		return bytes, nil
	}

	return base64.URLEncoding.DecodeString(buf.String())
}

// Returns a key of length "length" given an hkdf function
func getHKDFKey(hkdf io.Reader, length int) ([]byte, error) {
	key := make([]byte, length)
	n, err := io.ReadFull(hkdf, key)
	if n != len(key) || err != nil {
		return key, err
	}

	return key, nil
}

func pad(payload *bytes.Buffer, maxPadLen int) error {
	payloadLen := payload.Len()
	if payloadLen > maxPadLen {
		return ErrMaxPadExceeded
	}

	padLen := maxPadLen - payloadLen

	padding := make([]byte, padLen)
	payload.Write(padding)

	return nil
}
